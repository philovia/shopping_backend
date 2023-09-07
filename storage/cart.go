package storage

import (
	// "context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/mreym/go-fiber-postgres/models"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantupdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantBuyCart        = errors.New("cannot update the purchase")
)

func AddProductToCart(db *gorm.DB, productID uint, userID uint) error {
	var product models.ProductUser
	if err := db.First(&product, productID).Error; err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var user models.Users
	if err := db.First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	user.UserCart = append(user.UserCart, product)

	if err := db.Save(&user).Error; err != nil {
		log.Println(err)
		return ErrCantupdateUser
	}

	return nil
}

func RemoveCartItem(db *gorm.DB, productID uint, userID uint) error {
	var user models.Users
	if err := db.First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product models.ProductUser
	if err := db.First(&product, productID).Error; err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	for i, item := range user.UserCart {
		if item.Product_ID == productID {
			user.UserCart = append(user.UserCart[:i], user.UserCart[i+1:]...)
			break
		}
	}

	if err := db.Save(&user).Error; err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil
}

func BuyItemFromCart(db *gorm.DB, userID uint) error {
	var user models.Users
	if err := db.Preload("UserCart").First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var totalPrice int
	for _, item := range user.UserCart {
		totalPrice += item.Price
	}

	var order models.Order
	order.Ordered_At = time.Now()
	order.Order_Cart = user.UserCart
	order.Price = totalPrice
	order.Payment_Method.Digital = false
	order.Payment_Method.COD = true

	user.Order_Status = append(user.Order_Status, order)
	user.UserCart = nil

	if err := db.Save(&user).Error; err != nil {
		log.Println(err)
		return ErrCantBuyCart
	}

	return nil
}

func InstantBuyer(db *gorm.DB, productID uint, userID uint) error {
	var user models.Users
	if err := db.First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product models.ProductUser
	if err := db.First(&product, productID).Error; err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var order models.Order
	order.Ordered_At = time.Now()
	order.Order_Cart = []models.ProductUser{product}
	order.Price = product.Price
	order.Payment_Method.Digital = false
	order.Payment_Method.COD = true

	user.Order_Status = append(user.Order_Status, order)

	if err := db.Save(&user).Error; err != nil {
		log.Println(err)
		return ErrCantBuyCart
	}

	return nil
}
