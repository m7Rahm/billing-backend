package network

import (
	"billing/backend/database/repository"
	"billing/backend/models"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConnection struct {
	session *ssh.Session
	client  *ssh.Client
}
type NetworkServiceInterface interface {
	GetSwitchesDetailed() ([]models.Switch, error)
	GetSwitches(map[string]interface{}) ([]map[string]interface{}, error)
	GetSwitchInfo(context.Context, string) (*models.SwitchInfo, error)
	GetSwitchList(name string, vendor string, mac string, ip string) ([]map[string]interface{}, error)
	GetVendors() ([]models.Vendor, error)
	GetSwitchDetails(string) (*map[string]interface{}, error)
	AddNewSwitch(*models.Switch) error
	SetPortStatus(string, string, string) error
}
type NetworkService struct {
	netRepo repository.NetRepoInterface
}

func NewNetworkService(netRepo repository.NetRepoInterface) NetworkServiceInterface {
	return &NetworkService{netRepo: netRepo}
}

func (ns *NetworkService) AddNewSwitch(switchData *models.Switch) error {
	err := ns.netRepo.AddNewSwitch(switchData)
	return err
}
func (ns *NetworkService) SetPortStatus(switchId string, portId string, state string) error {
	// switches, err := ns.netRepo.GetSwitchesDetailed()
	if switchId == "" {
		return fmt.Errorf("switch id is empty")
	}
	portState := ""
	if state == "1" {
		portState = "enable"
	} else if state == "0" {
		portState = "disable"
	} else {
		return fmt.Errorf("invalid port state")
	}
	switchIdInt, err := strconv.ParseInt(switchId, 10, 64)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	sshConnection, err := ns.connectToSwitch(int(switchIdInt), ctx)
	if err != nil {
		return err
	}
	session := sshConnection.session
	client := sshConnection.client
	defer client.Close()
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	w, err := session.StdinPipe()
	if err != nil {
		return err
	}
	if err = session.RequestPty("xterm", 180, 140, ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		return err
	}
	setPortState := fmt.Sprintf("interface ethernet %s %s", portId, portState)
	commands := []string{
		"debug destination session",
		"configure terminal",
		setPortState,
	}
	if err != nil {
		return err
	}
	err = session.Shell()
	if err != nil {
		return err
	}
	debugMessage := ""
	if portState == "enable" {
		debugMessage = fmt.Sprintf("ports: port %s is now on-line", portId)
	} else {
		debugMessage = fmt.Sprintf("ports: port %s is now off-line", portId)
	}
	w.Write([]byte("\n"))
	for _, command := range commands {
		w.Write([]byte(fmt.Sprintf("%s\n", command)))
	}
	buffer := make([]byte, 512)
	str := ""
	n := 1
	hasExitRequested := false
	for n > 0 {
		n, err = stdout.Read(buffer)
		if err != nil {
			log.Println(err, n)
		}
		if strings.Contains(string(buffer[:n]), debugMessage) && !hasExitRequested {
			w.Write([]byte("exit\n"))
			w.Write([]byte("exit\n"))
			w.Write([]byte("exit\n"))
			hasExitRequested = true
		} else if strings.Contains(string(buffer[:n]), "Do you want to log out (y/n)?") {
			w.Write([]byte("y"))
			break
		}
		str += string(buffer[:n])
	}
	return nil
}
func (ns *NetworkService) GetSwitchDetails(switchId string) (*map[string]interface{}, error) {
	return ns.netRepo.GetSwitchDetails(switchId)
}
func (ns *NetworkService) GetVendors() ([]models.Vendor, error) {
	vendors, err := ns.netRepo.GetVendors()
	if err != nil {
		return []models.Vendor{}, err
	}
	return vendors, nil
}
func (ns *NetworkService) GetSwitches(query map[string]interface{}) ([]map[string]interface{}, error) {
	switches := ns.netRepo.GetSwitches(query)
	return switches, nil
}
func (ns *NetworkService) GetSwitchList(name string, vendor string, mac string, ip string) ([]map[string]interface{}, error) {
	switches, err := ns.netRepo.GetSwitchList(name, vendor, mac, ip)
	return switches, err
}
func (ns *NetworkService) GetSwitchInfo(ctx context.Context, switchId string) (*models.SwitchInfo, error) {
	switchIdInt, err := strconv.ParseInt(switchId, 10, 64)
	if err != nil {
		return &models.SwitchInfo{}, err
	}
	sshConnection, err := ns.connectToSwitch(int(switchIdInt), ctx)
	if err != nil {
		return &models.SwitchInfo{}, err
	}
	session := sshConnection.session
	client := sshConnection.client
	defer client.Close()
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return &models.SwitchInfo{}, err
	}

	w, err := session.StdinPipe()
	if err != nil {
		return &models.SwitchInfo{}, err
	}
	if err = session.RequestPty("xterm", 180, 140, ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		return &models.SwitchInfo{}, err
	}
	commands := []string{
		"no page",
		"show lldp info remote-device",
		"show interface status",
		"exit",
		"exit",
	}
	if err != nil {
		return &models.SwitchInfo{}, err
	}
	err = session.Shell()
	if err != nil {
		return &models.SwitchInfo{}, err
	}
	w.Write([]byte("\n"))
	for _, command := range commands {
		w.Write([]byte(fmt.Sprintf("%s\n", command)))
	}
	buffer := make([]byte, 512)
	str := ""
	n := 1
	for n > 0 {
		n, err = stdout.Read(buffer)
		if err != nil {
			log.Println(err, n)
		}
		if strings.Contains(string(buffer[:n]), "Do you want to log out (y/n)?") {
			w.Write([]byte("y"))
			break
		}
		str += string(buffer[:n])
	}
	commandResults := strings.Split(str, "HP-2530-48G#")
	if len(commandResults) < 4 {
		return &models.SwitchInfo{}, fmt.Errorf("command results not found")
	}
	commandLines := make([][]string, len(commandResults))
	for i := 2; i < 4; i++ {
		lines := strings.Split(commandResults[i], "\n")
		commandLines[i] = append(commandLines[i], lines...)
	}
	remoteDeviceInfo := commandLines[2][5:]
	interfaceStatus := commandLines[3][2:]
	var remoteDevices []models.RemoteDevice
	re := regexp.MustCompile(`^\x20*(?P<port>\d+)\x20*\| (?P<mac>([a-z0-9/x20]{2} ){5}[a-z0-9/x20]{2})  (?P<portId>.{18}) (?P<portDescr>.{9}) (?P<sysName>.*)$`)
	for i := 0; i < len(remoteDeviceInfo); i++ {
		if len(remoteDeviceInfo[i]) < 62 {
			continue
		}
		result := make(map[string]string)
		match := re.FindStringSubmatch(remoteDeviceInfo[i])
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		localPort, err := strconv.ParseInt(strings.Trim(result["port"], " "), 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		mac := strings.Trim(result["mac"], " ")
		portId := strings.Trim(result["portId"], " ")
		portDesc := strings.Trim(result["portDescr"], " ")
		sysName := strings.Trim(result["sysName"], " ")
		device := models.RemoteDevice{
			LocalPort: int(localPort),
			Mac:       mac,
			PortId:    portId,
			PortDesc:  portDesc,
			SysName:   sysName,
		}
		remoteDevices = append(remoteDevices, device)
	}
	var interfaces []models.Interface
	for i := 0; i < len(interfaceStatus); i++ {
		if len(interfaceStatus[i]) < 60 {
			continue
		}
		port := strings.Trim(interfaceStatus[i][2:10], " ")
		name := strings.Trim(interfaceStatus[i][11:21], " ")
		status := strings.Trim(interfaceStatus[i][22:29], " ")
		configMode := strings.Trim(interfaceStatus[i][30:43], " ")
		speed := strings.Trim(interfaceStatus[i][44:52], " ")
		typ := strings.Trim(interfaceStatus[i][53:63], " ")
		tagged := strings.Trim(interfaceStatus[i][64:70], " ")
		untagged := strings.Trim(interfaceStatus[i][71:79], " ")
		intf := models.Interface{
			Port:       port,
			Name:       name,
			Status:     status,
			ConfigMode: configMode,
			Speed:      speed,
			Type:       typ,
			Tagged:     tagged,
			Untagged:   untagged,
		}
		interfaces = append(interfaces, intf)
	}
	return &models.SwitchInfo{
		RemoteDevices: remoteDevices,
		Interfaces:    interfaces,
	}, nil
}
func (c *NetworkService) GetSwitchesDetailed() ([]models.Switch, error) {
	// switches, err := c.netRepo.GetSwitchesDetailed()
	return []models.Switch{}, nil
}
