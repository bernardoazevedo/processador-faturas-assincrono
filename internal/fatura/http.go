package fatura

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func HttpProcessaFaturas(c *gin.Context) {
	var faturas []Fatura

	err := c.BindJSON(&faturas)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error parsing json, verify your request: " + err.Error()})
		return
	}

	err = ProcessaFaturas(faturas)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.IndentedJSON(http.StatusOK, gin.H{"faturas": faturas})
}

func HttpListaFaturas(c *gin.Context) {
	faturas, err := ListaFaturas()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"faturas": faturas})
}