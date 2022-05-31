package models

import (
	"gorm.io/gorm"
)

type GoModel gorm.Model
type RemoteDevice struct {
	LocalPort int    `json:"local_port"`
	PortId    string `json:"port_id"`
	Mac       string `json:"mac"`
	PortDesc  string `json:"port_desc"`
	SysName   string `json:"sys_name"`
}

type SwitchInfo struct {
	RemoteDevices []RemoteDevice `json:"remote_devices"`
	Interfaces    []Interface    `json:"interfaces"`
}
type Vendor struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}
type Interface struct {
	Port       string `json:"port"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	ConfigMode string `json:"config_mode"`
	Speed      string `json:"speed"`
	Type       string `json:"type"`
	Tagged     string `json:"tagged"`
	Untagged   string `json:"untagged"`
}
type Switch struct {
	GoModel
	Mac          string   `json:"mac"`
	Name         string   `json:"name"`
	IP           string   `json:"ip"`
	Serial       string   `json:"serial"`
	Username     string   `json:"username"`
	VendorID     int      `json:"vendor_id"`
	Password     string   `json:"password"`
	Model        string   `json:"model"`
	PortCount    int      `json:"port_count"`
	CityID       int      `json:"city_id"`
	City         *Address `json:"city" gorm:"foreignkey:CityID"`
	SettlementID *int     `json:"settlement_id"`
	Settlement   *Address `json:"settlement" gorm:"foreignkey:SettlementID"`
	DistrictID   *int     `json:"district_id"`
	District     *Address `json:"district" gorm:"foreignkey:DistrictID"`
	BuildingID   int      `json:"building_id"`
	Building     *Address `json:"building" gorm:"foreignkey:BuildingID"`
	StreetId     int      `json:"street_id"`
	Street       *Address `json:"street" gorm:"foreignkey:StreetId"`
}

type Address struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	AddressType int    `json:"address_type"`
}
type SwitchPort struct {
	ID  int    `json:"id" gorm:"primary_key"`
	Mac string `json:"mac"`
}
type AddressType struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
