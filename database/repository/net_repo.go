package repository

import (
	"billing/backend/models"
	"context"

	// "context"

	"gorm.io/gorm"
)

type NetRepoInterface interface {
	GetSwitches(map[string]interface{}) ([]map[string]interface{}, error)
	GetSwitch(int) (models.Switch, error)
	GetSwitchDetails(switchId string) (*map[string]interface{}, error)
	GetSSHConfig(context.Context, int) (map[string]interface{}, error)
	GetSwitchList(name string, vendor string, mac string, ip string) ([]map[string]interface{}, error)
	GetVendors() ([]models.Vendor, error)
	AddNewSwitch(*models.Switch) error
}

type NetRepo struct {
	db *gorm.DB
}

func (nr *NetRepo) GetSwitches(query map[string]interface{}) ([]map[string]interface{}, error) {
	var switches []map[string]interface{}
	err := nr.db.Table("switches").Select("id, mac, ip, name, port_count, city_id, building_id, street_id").Where(query).Find(&switches).Error
	return switches, err
}
func (nr *NetRepo) AddNewSwitch(switchInfo *models.Switch) error {
	return nr.db.Create(switchInfo).Error
}
func (nr *NetRepo) GetVendors() ([]models.Vendor, error) {
	var vendors []models.Vendor
	err := nr.db.Model(&models.Vendor{}).Scan(&vendors).Error
	return vendors, err
}
func (nr *NetRepo) GetSwitchList(name string, vendor string, mac string, ip string) ([]map[string]interface{}, error) {
	var switchList []map[string]interface{}
	err := nr.db.
		Table("switches").
		Select("switches.id, switches.mac, switches.ip, vendors.name as vendor_name, model, switches.name, port_count").
		Joins("left join vendors on switches.vendor_id = vendors.id").
		Where("switches.name like ? AND switches.ip like ? AND switches.mac like ? AND vendors.name like ? ", "%"+name+"%", "%"+ip+"%", "%"+mac+"%", "%"+vendor+"%").
		Scan(&switchList).
		Error
	return switchList, err
}
func (n *NetRepo) GetSwitchDetails(switchId string) (*map[string]interface{}, error) {
	var switchDetailes map[string]interface{}
	err := n.db.Model(&models.Switch{}).
		Select("switches.id, switches.mac, switches.ip, model, switches.name, port_count, serial, city_id, street_id, building_id").
		Where("id = ?", switchId).Scan(&switchDetailes).
		Error
	if err != nil {
		return nil, err
	}
	return &switchDetailes, nil
}
func (n *NetRepo) GetSwitch(id int) (models.Switch, error) {
	return models.Switch{}, nil
}
func (n *NetRepo) GetSSHConfig(ctx context.Context, id int) (map[string]interface{}, error) {
	var switchSSHInfo map[string]interface{}
	err := n.db.Table("switches").
		WithContext(ctx).
		Select("ip, username, password").
		Where("id = ?", id).
		Scan(&switchSSHInfo).Error
	return switchSSHInfo, err
}
func NewNetRepo(db *gorm.DB) NetRepoInterface {
	return &NetRepo{db: db}
}
