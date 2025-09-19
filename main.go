package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", PingPong)
	r.GET("/newping", NewPong)
	if err := r.Run(); err != nil {
		panic(err)
	}
}

func PingPong(c *gin.Context) {
	c.String(200, "pong")
}

func NewPong(c *gin.Context) {
	c.String(200, "new pong")
}
