package models

type ProductLocation struct {
	ProductID   uint `gorm:"primaryKey;column:product_id" json:"productId"`
	LocationID  uint `gorm:"primaryKey;column:location_id" json:"locationId"`
	IsAvailable bool `gorm:"column:is_available;default:true" json:"isAvailable"`

	Product  *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product,omitempty"`
	Location *Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
}
