package models

type Resto struct {
  ID     uint   `json:"id" gorm:"primary_key"`
  Nama  string `json:"nama"`
  Alamat string `json:"alamat"`
  Lat float64 `json:"lat"`
  Lng float64 `json:"lng"`
  Radius int `json:"radius"` // valid radius
}

type CreateRestoInput struct {
	Nama  string `json:"nama" binding:"required"`
	Alamat string `json:"alamat" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
	Radius int `json:"radius" binding:"required"`
}

type UpdateRestoInput struct {
	// Nama  string `json:"nama" binding:"required"`
	// Alamat string `json:"alamat" binding:"required"`
	// Lat float64 `json:"lat" binding:"required"`
	// Lng float64 `json:"lng" binding:"required"`
	Radius int `json:"radius" binding:"required"`
}