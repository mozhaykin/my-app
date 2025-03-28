package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Profile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	router := gin.Default()

	router.GET("/amozhaykin/my-app/hello", helloGET)
	router.GET("/amozhaykin/my-app/profile", profileGET)
	router.POST("/amozhaykin/my-app/profile", profilePOST)

	log.Info().Msg("Starting server on :8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func helloGET(c *gin.Context) {
	log.Debug().Msg("Handling GET /hello")

	c.JSON(200, "Hello!")
}

func profileGET(c *gin.Context) {
	log.Debug().Msg("Handling GET /profile")

	Profile := Profile{
		Name: "Alice",
		Age:  30,
	}

	c.JSON(200, Profile)
}

func profilePOST(c *gin.Context) {
	log.Debug().Msg("Handling POST /profile")

	var newProfile Profile

	err := c.BindJSON(&newProfile)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	c.Status(201)
}
