package models

type Client struct {
	GoModel
	FullName         string `json:"full_name"`
	SwitchPort       int    `json:"switch_port"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	PassportData     string `json:"passport_data"`
	PassportFilePath string `json:"passport_file_path"`
	ContractId       int    `json:"contract_id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Ip               string `json:"ip"`
	IpV6             string `json:"ipv6"`
	Mac              string `json:"mac"`
	SwitchId         int    `json:"switch_id"`
}
