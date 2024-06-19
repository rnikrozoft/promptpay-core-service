package repository

import (
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *repository { return &repository{db: db} }
