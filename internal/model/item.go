package model

type Item struct {
	ID         uint64 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	ItemID     string `gorm:"column:item_id;type:varchar(50);not null" json:"itemID"`
	ItemName   string `gorm:"column:item_name;type:varchar(50);not null" json:"itemName"`
	ItemCname  string `gorm:"column:item_cname;type:varchar(50)" json:"itemCname"`
	ItemDesc   string `gorm:"column:item_desc;type:text" json:"itemDesc"`
	Hp         int    `gorm:"column:hp;type:int(11)" json:"hp"`
	Mp         int    `gorm:"column:mp;type:int(11)" json:"mp"`
	Attack     int    `gorm:"column:attack;type:int(11)" json:"attack"`
	Defence    int    `gorm:"column:defence;type:int(11)" json:"defence"`
	Dodge      int    `gorm:"column:dodge;type:int(11)" json:"dodge"`
	Str        int    `gorm:"column:str;type:int(11)" json:"str"`
	Cor        int    `gorm:"column:cor;type:int(11)" json:"cor"`
	Inte       int    `gorm:"column:inte;type:int(11)" json:"inte"`
	Dex        int    `gorm:"column:dex;type:int(11)" json:"dex"`
	Con        int    `gorm:"column:con;type:int(11)" json:"con"`
	Kar        int    `gorm:"column:kar;type:int(11)" json:"kar"`
	Classifier string `gorm:"column:classifier;type:varchar(1)" json:"classifier"`
}

// TableName table name
func (m *Item) TableName() string {
	return "item"
}

// ItemColumnNames Whitelist for custom query fields to prevent sql injection attacks
var ItemColumnNames = map[string]bool{
	"id":         true,
	"item_id":    true,
	"item_name":  true,
	"item_cname": true,
	"item_desc":  true,
	"hp":         true,
	"mp":         true,
	"attack":     true,
	"defence":    true,
	"dodge":      true,
	"str":        true,
	"cor":        true,
	"inte":       true,
	"dex":        true,
	"con":        true,
	"kar":        true,
	"classifier": true,
}
