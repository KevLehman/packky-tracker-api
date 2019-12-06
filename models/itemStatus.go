package models

import "github.com/jinzhu/gorm"

import "fmt"

// ItemStatus defines the valid statuses an item can have
type ItemStatus struct {
	gorm.Model
	Name        string
	Description string
	ImageURI    string
	Items       []Item
}

// GetInitialStatus returns the default initial status for a product
func (i ItemStatus) GetInitialStatus(db *gorm.DB) (*ItemStatus, error) {
	status := &ItemStatus{}
	err := db.Model(&i).Where("name = ?", "En Espera").First(&status).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

// GetStatusByName returns the first ItemStatus instance matching status on DB
func (i ItemStatus) GetStatusByName(db *gorm.DB, status string) (*ItemStatus, error) {
	stat := &ItemStatus{}
	err := db.Model(&i).Where("name ilike ?", status).First(&stat).Error
	if err != nil {
		return nil, err
	}
	return stat, nil
}

// Seed populates DB with default Order statuses
func (i ItemStatus) Seed(db *gorm.DB) {
	tx := db.Begin()
	statuses := []ItemStatus{
		ItemStatus{
			Name:        "En Espera",
			Description: "Default starting status",
			ImageURI:    "https://snipboard.io/FOt8CI.jpg",
		},
		ItemStatus{
			Name:        "En Camino",
			Description: "Order was placed & is on its way to be delivered",
			ImageURI:    "https://snipboard.io/EcrhWD.jpg",
		},
		ItemStatus{
			Name:        "Entregada",
			Description: "Order was delivered on time to ",
			ImageURI:    "https://snipboard.io/rVxOWw.jpg",
		},
		ItemStatus{
			Name:        "Cancelada",
			Description: "Order was cancelled, either by User or Client",
			ImageURI:    "https://snipboard.io/d6qFIG.jpg",
		},
		ItemStatus{
			Name:        "No entregada",
			Description: "Order was not delivered because of inability to contact receiver or tech problems",
			ImageURI:    "https://snipboard.io/d6qFIG.jpg",
		},
	}

	for _, status := range statuses {
		fmt.Println(status)
		if err := tx.Save(&status).Error; err != nil {
			tx.Rollback()
			break
		}
	}

	tx.Commit()
}
