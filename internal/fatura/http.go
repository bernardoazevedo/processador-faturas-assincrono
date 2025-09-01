package fatura

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HttpProcessList(c *gin.Context) {
	var faturas []Fatura

	err := c.BindJSON(&faturas)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error parsing json, verify your request: " + err.Error()})
		return
	}

	err = ProcessList(faturas)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.IndentedJSON(http.StatusOK, gin.H{"faturas": faturas})
}

func HttpList(c *gin.Context) {
	faturas, err := List()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"faturas": faturas})
}