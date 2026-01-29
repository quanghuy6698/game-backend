package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"game-backend/src/module/tetGame"
)

func main() {
	route := gin.Default()

	route.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := tetGame.InitDB(); err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer tetGame.CloseDB()

	tetGameApi := route.Group("/tet-game")
	{
		tetGameApi.POST("/create-player", tetGame.CreatePlayer)
		tetGameApi.GET("/top-100", tetGame.GetTop100)
		tetGameApi.GET("/save-score", tetGame.SaveScore)
	}

	log.Println("Server running on :8080")
	route.Run(":8080")
}
