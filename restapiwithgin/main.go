package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users   = make(map[int]User)
	nextID  = 1
	usersMu sync.Mutex
)

func main() {
	router := gin.Default()

	// POST /users
	router.POST("/users", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		usersMu.Lock()
		newUser.ID = nextID
		nextID++
		users[newUser.ID] = newUser
		usersMu.Unlock()

		c.JSON(http.StatusCreated, newUser)
	})

	// GET /users
	router.GET("/users", func(c *gin.Context) {
		usersMu.Lock()
		userList := make([]User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}
		usersMu.Unlock()

		c.JSON(http.StatusOK, userList)
	})

	// Start the server
	router.Run(":8080") // Default is localhost:8080
}
