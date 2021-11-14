package models

import (
	"gorm.io/gorm"
)

type LogScan struct {
	gorm.Model
	RestoID int
  	Resto   Resto
	UserName string `json:"username"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}  


type CreateLogScanInput struct {
	RestoID int `json:"restoId" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
} 

func (b *LogScan) TableName() string {
	return "log_scan"
}