package main

import (
	"github.com/gin-gonic/gin"
)

var Users map[string]*User

type User struct {
	Id    string
	Email string
	Name  string
}

func GetCertificatesHandler(c *gin.Context) {

}

func GetUserOrError(c *gin.Context) *User {
	Session := c.Request.Header.Get("Authorization")
	if Session == "" {
		c.AbortWithStatus(401)
		return nil
	}
	if UserID, Exists := Sessions[Session]; Exists {
		return Users[UserID]
	} else {
		c.AbortWithStatus(401)
		return nil
	}
}
