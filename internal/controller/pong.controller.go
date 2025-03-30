package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PongController struct{}

func NewPongController() *PongController {
	return &PongController{}
}

func (pc *PongController) Pong(c *gin.Context) {
	fmt.Println("--> My handler")
	name := c.DefaultQuery("name", "quangduong")
	c.JSON(http.StatusOK, gin.H{
		"message": "ping....pong" + name,
		"users":   []string{"qduong", "defnotqduong"},
	})
}
