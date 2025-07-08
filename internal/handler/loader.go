package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadVaultData(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"success": true})
}
