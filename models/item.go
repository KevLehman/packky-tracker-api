package models

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

// Item represents an order placed on API
type Item struct {
	gorm.Model
	Name        string
	Description string
	ItemID      string `gorm:"unique; not null"`
	StatusID    int
	Status      ItemStatus
}

// BeforeSave performs data validation before persisting item on DB
func (i *Item) BeforeSave() error {
	if !i.Validate() {
		return errors.New("Validation error while creating Item instance")
	}
	i.ItemID = strings.ToUpper(i.ItemID)

	return nil
}

// Validate checks if model's data is valid
func (i *Item) Validate() bool {
	itemRegexp, err := regexp.Compile("^(?i)serv-[0-9]{5}$")
	if err != nil {
		log.Println(err)
		return false
	}

	if !itemRegexp.Match([]byte(i.ItemID)) {
		return false
	}
	if i.Name == "" {
		return false
	}

	return true
}

// GetItemByOrder returns the first instance of Item with matching order
func (i Item) GetItemByOrder(db *gorm.DB, order string) (*Item, error) {
	logger := &Log{}
	out := &Item{}
	err := db.Model(&i).Preload("Status").Where("item_id = ?", strings.ToUpper(order)).First(&out).Error

	if err != nil {
		logger.TrackAction(db, "find-or-create-error", err.Error())
		return nil, err
	}
	return out, nil
}

// UpdateStatus changes the status of item
func (i *Item) UpdateStatus(db *gorm.DB, status string) (*Item, error) {
	logger := &Log{}
	tx := db.Begin()

	st := ItemStatus{}

	err := tx.Find(&st, "name ilike ?", status).Error
	if err != nil {
		logger.TrackAction(db, "update-item-status-error", err.Error())
		tx.Rollback()
		return nil, err
	}

	i.Status = st
	if err := tx.Save(&i).Error; err != nil {
		logger.TrackAction(db, "update-item-status-error", err.Error())
		tx.Rollback()
		return nil, err
	}

	logger.TrackAction(tx, "update-item-status", fmt.Sprintf("Order %s to status %s", i.ItemID, status))
	tx.Commit()

	return i, nil
}

// FindOrCreate returns the instance of order of database or creates it on DB
func (i Item) FindOrCreate(db *gorm.DB, order string) (*Item, error) {
	logger := &Log{}
	foundOrder, err := Item{}.GetItemByOrder(db, order)

	logger.TrackAction(db, "find-or-create-order", fmt.Sprintf("Order: %s", order))

	if err != nil {
		logger.TrackAction(db, "find-or-create-error", err.Error())
		if gorm.IsRecordNotFoundError(err) {
			newItem := &Item{
				Name:        fmt.Sprintf("Order %s", order),
				ItemID:      strings.ToUpper(order),
				Description: "Dummy description | Description on packages not supported yet",
			}
			err = db.Save(&newItem).Error
			if err != nil {
				return nil, err
			}
			return newItem, nil
		}
		return nil, err
	}

	return foundOrder, nil
}

// Update will save the instance of user on the database
func (i *Item) Update(db *gorm.DB) error {
	l := Log{}
	err := db.Save(&i).Error
	if err != nil {
		l.TrackAction(db, "update-item-error", err.Error())
		return err
	}
	return nil
}

// Seed will populate DB with dummy item data
func (i Item) Seed(db *gorm.DB) {
	tx := db.Begin()
	initialStatus, _ := ItemStatus{}.GetStatusByName(tx, "En Espera")
	otherStatus, _ := ItemStatus{}.GetStatusByName(tx, "En Camino")
	items := []Item{
		Item{
			Name:        "Order XXXX",
			ItemID:      "SERV-12345",
			Description: "A random order",
			Status:      *initialStatus,
		},
		Item{
			Name:        "Order XYZAY",
			ItemID:      "SERV-65331",
			Description: "A random order II",
			Status:      *otherStatus,
		},
	}

	for _, item := range items {
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
		}
	}

	tx.Commit()
}
