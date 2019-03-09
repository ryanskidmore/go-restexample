package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

var Users map[string]*User

type User struct {
	Id    string
	Email string
	Name  string
}

func GetCertificatesHandler(c *gin.Context) {
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	RequestedID := c.Param("userId")
	if RequestedID == "" {
		c.AbortWithStatus(400)
		log.Println("GET /users/:userId/certificates: Failed to get users certificates (missing userId parameter)")
		return
	}
	UserCertificates := []*Certificate{}
	if User.Id == RequestedID {
		for _, Certificate := range Certificates {
			if Certificate.OwnerId == User.Id {
				UserCertificates = append(UserCertificates, Certificate)
			}
		}
	} else {
		c.AbortWithStatus(401)
		log.Println("GET /users/:userId/certificates: Failed to get users certificates (requesting another users certificates)")
		return
	}
	c.JSON(200, map[string]interface{}{
		"status":       "success",
		"certificates": UserCertificates,
	})
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
