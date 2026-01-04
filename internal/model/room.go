package model

type Room struct {
	ID    string `gorm:"column:id;type:varchar(50);primary_key" json:"id"`
	Title string `gorm:"column:title;type:varchar(30);not null" json:"title"`
	Desc  string `gorm:"column:desc;type:text" json:"desc"`
	Way   string `gorm:"column:way;type:varchar(30)" json:"way"`
	Mobs  string `gorm:"column:mobs;type:varchar(256)" json:"mobs"`
}

// TableName table name
func (m *Room) TableName() string {
	return "room"
}

// RoomColumnNames Whitelist for custom query fields to prevent sql injection attacks
var RoomColumnNames = map[string]bool{
	"id":    true,
	"title": true,
	"desc":  true,
	"way":   true,
	"mobs":  true,
}
