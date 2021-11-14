package controllers

import (
	m "RestApi/models"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// jwt "github.com/appleboy/gin-jwt/v2"
)

// GET /books
// Get all books
func FindLogScan(c *gin.Context) {
	pagination := GeneratePaginationFromRequest(c)
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuider := m.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)

	var models []m.LogScan
	queryBuider.Model(&m.LogScan{}).Preload("Resto").Find(&models)
	// m.DB.Preload("Resto").Find(&models)
	var count int64
	m.DB.Model(&m.LogScan{}).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": models,
			"total": count,
		},
	})
}