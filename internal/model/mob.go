package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type Mob struct {
	ID         uint64          `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	MobID      string          `gorm:"column:mob_id;type:varchar(50);not null" json:"mobID"`
	MobName    string          `gorm:"column:mob_name;type:varchar(50);not null" json:"mobName"`
	MobCname   string          `gorm:"column:mob_cname;type:varchar(50);not null" json:"mobCname"`
	MobDesc    string          `gorm:"column:mob_desc;type:text" json:"mobDesc"`
	Attackable *sgorm.TinyBool `gorm:"column:attackable;type:tinyint(1)" json:"attackable"`
	Hp         int             `gorm:"column:hp;type:int(11);default:100;not null" json:"hp"`
	Mp         int             `gorm:"column:mp;type:int(11);default:100;not null" json:"mp"`
	Attack     int             `gorm:"column:attack;type:int(11);default:1;not null" json:"attack"`
	Defence    int             `gorm:"column:defence;type:int(11);default:1;not null" json:"defence"`
	Dodge      int             `gorm:"column:dodge;type:int(11);default:1;not null" json:"dodge"`
}

// TableName table name
func (m *Mob) TableName() string {
	return "mob"
}

// MobColumnNames Whitelist for custom query fields to prevent sql injection attacks
var MobColumnNames = map[string]bool{
	"id":         true,
	"mob_id":     true,
	"mob_name":   true,
	"mob_cname":  true,
	"mob_desc":   true,
	"attackable": true,
	"hp":         true,
	"mp":         true,
	"attack":     true,
	"defence":    true,
	"dodge":      true,
}
