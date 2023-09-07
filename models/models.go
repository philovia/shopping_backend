package models

import (
	"time"
	// "github.com/mreym/go-fiber-postgres/address"
)

type Users struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	First_Name      string        `gorm:"type:varchar(30);not null" json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       string        `gorm:"type:varchar(30);not null" json:"last_name" validate:"required,min=2,max=30"`
	Password        string        `gorm:"not null" json:"-" validate:"required,min=6"`
	Email           string        `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
	Phone           string        `gorm:"type:varchar(15);not null" json:"phone" validate:"required"`
	Token           string        `gorm:"-" json:"token"`
	Refresh_Token   string        `gorm:"-" json:"refresh_token"`
	Created_At      time.Time     `json:"created_at"`
	Updated_At      time.Time     `json:"updated_at"`
	User_ID         int           `json:"user_id"`
	UserCart        []ProductUser `gorm:"foreignKey:UserID" json:"usercart"`
	Address_Details []Address     `gorm:"foreignKey:UserID" json:"address"`
	Order_Status    []Order       `gorm:"foreignKey:UserID" json:"order"`
}

type Product struct {
	Product_ID   uint    `gorm:"primaryKey" json:"product_id"`
	Product_Name string  `gorm:"type:varchar(100)" json:"product_name"`
	Price        int     `json:"price"`
	Rating       *uint   `json:"rating"`
	Image        *string `gorm:"type:varchar(255)" json:"image"`
}

type ProductUser struct {
	Product_ID  uint    `gorm:"primaryKey" json:"product_user_id"`
	UserID      uint    `gorm:"index" json:"-"`
	ProductName string  `gorm:"type:varchar(100)" json:"product_name"`
	Price       int     `json:"price"`
	Rating      *uint   `json:"rating"`
	Image       *string `gorm:"type:varchar(255)" json:"image"`
}

type Address struct {
	Address_ID uint   `gorm:"primaryKey" json:"address_id"`
	UserID     uint   `gorm:"index" json:"-"`
	House      string `gorm:"type:varchar(255)" json:"house_name"`
	Street     string `gorm:"type:varchar(255)" json:"street_name"`
	City       string `gorm:"type:varchar(100)" json:"city_name"`
	Pincode    string `gorm:"type:varchar(10)" json:"pin_code"`
}

type Order struct {
	Order_ID       uint          `gorm:"primaryKey" json:"order_id"`
	Order_Cart     []ProductUser `json:"order_list" bson:"order_list"`
	UserID         uint          `gorm:"index" json:"-"`
	Ordered_At     time.Time     `json:"ordered_at"`
	Price          int           `json:"total_price"`
	Discount       *int          `json:"discount"`
	Payment_Method Payment       `gorm:"embedded" json:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital"`
	COD     bool `json:"cod"`
}
