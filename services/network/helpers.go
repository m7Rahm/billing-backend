package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func (ns *NetworkService) connectToSwitch(switchId int, ctx context.Context) (*SSHConnection, error) {
	switchInfo, err := ns.netRepo.GetSSHConfig(ctx, switchId)
	if err != nil {
		return nil, err
	}
	ip, ok := switchInfo["ip"]
	if !ok {
		return nil, fmt.Errorf("ip not found")
	}
	username, ok := switchInfo["username"]
	if !ok {
		return nil, fmt.Errorf("username not found")
	}
	password, ok := switchInfo["password"]
	if !ok {
		return nil, fmt.Errorf("password not found")
	}
	d := net.Dialer{Timeout: time.Minute}
	cntx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := d.DialContext(cntx, "tcp", fmt.Sprintf("%s:22", ip))
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, fmt.Sprintf("%s:22", ip), &ssh.ClientConfig{
		User: fmt.Sprintf("%s", username),
		Auth: []ssh.AuthMethod{
			ssh.Password(fmt.Sprintf("%s", password)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	client := ssh.NewClient(c, chans, reqs)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return &SSHConnection{
		client:  client,
		session: session,
	}, nil
}
