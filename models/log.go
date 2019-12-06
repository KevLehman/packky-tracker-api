package models

import "github.com/jinzhu/gorm"

import "strings"

// Log saves a history of actions performed on api
type Log struct {
	gorm.Model
	Action   string
	IsError  bool
	Metadata string
}

// TrackAction stores a permanent log of an action performed on API
func (l *Log) TrackAction(db *gorm.DB, action string, metadata string) {
	l.Action = action
	l.Metadata = metadata
	l.IsError = strings.Contains(action, "error")

	db.Save(&l)
}
