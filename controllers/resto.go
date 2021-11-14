package controllers

import (
	m "RestApi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umahmood/haversine"
	"gorm.io/gorm"
	// jwt "github.com/appleboy/gin-jwt/v2"

)

// GET /books
// Get all books
// func FindResto(c *gin.Context) {
// 	var models []m.Resto
// 	m.DB.Find(&models)
// 	c.JSON(http.StatusOK, gin.H{"data": models})
// }

func FindResto(c *gin.Context) {
	pagination := GeneratePaginationFromRequest(c)
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuider := m.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)

	var models []m.Resto
	queryBuider.Model(&m.Resto{}).Find(&models)
	// m.DB.Preload("Resto").Find(&models)
	var count int64
	m.DB.Model(&m.Resto{}).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": models,
			"total": count,
		},
	})
}

func FindRestoById(c *gin.Context) {  // Get model if exist
	var model m.Resto
  
	if err := m.DB.Where("id = ?", c.Param("id")).First(&model).Error; err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	  return
	}
  
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": model,
	})
  }

func CreateResto(c *gin.Context) {
	// Validate input
	var input m.CreateRestoInput
	if err := c.ShouldBindJSON(&input); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
  
	// Create book
	model := m.Resto{Nama: input.Nama, Alamat: input.Alamat, Lat: input.Lat, Lng: input.Lng, Radius: input.Radius}
	m.DB.Create(&model)
  
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": model,
	})
}


func UpdateResto(c *gin.Context) {
	// Validate input
	var model m.Resto
	if err := m.DB.Where("id = ?", c.Param("id")).First(&model).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := c.ShouldBindJSON(&model); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}

	m.DB.Updates(model)
  
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": model,
	})
}

var identityKey = "id"

func ScanQR(c *gin.Context) {
	// create
	var input m.CreateLogScanInput
	if err := c.ShouldBindJSON(&input); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}

	var resto m.Resto
	err := m.DB.First(&resto, input.RestoID).Error
	if (err == gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resto not found", "data": resto})
		return
	}
	restoCoor := haversine.Coord{Lat: resto.Lat, Lon: resto.Lng}  // Oxford, UK
    userCoor  := haversine.Coord{Lat: input.Lat, Lon: input.Lng}
	_, km := haversine.Distance(restoCoor, userCoor)

	user, _ := c.Get(identityKey)
	username := user.(*m.User).UserName
	
	if(km * 1000 <= float64(resto.Radius)) {
		model := m.LogScan{RestoID: input.RestoID, Lat: input.Lat,Lng: input.Lng, UserName: username}
		m.DB.Preload("Resto").Create(&model)
		c.JSON(http.StatusOK, gin.H{"data": model})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "radius not valid"})
	}
}