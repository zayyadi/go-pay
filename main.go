package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zayyadi/go-pay/services/logic"
)

func main() {
	var GIN_MODE string = "release"

	mode := GIN_MODE
	if mode == "" {
		mode = "debug" // default to debug mode
	}
	gin.SetMode(mode)
	router := gin.Default()
	router.Use(logic.CORSMiddleware())
	router.LoadHTMLGlob("templates/*")

	// Serve static files
	router.Static("/static", "./static")
	router.GET("/", logic.Home)
	router.POST("/payslip", logic.GetPayslip)
	router.Run(":8098")
}
