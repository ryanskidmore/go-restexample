package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

// Map of users
// User ID -> *User
var Users map[string]*User

// User object
type User struct {
	Id    string
	Email string
	Name  string
}

// Handler for getting a users certificates
func GetCertificatesHandler(c *gin.Context) {
	// Get the user from provided authorisation token or throw error
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	// Get user id from url parameters
	RequestedID := c.Param("userId")
	if RequestedID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: missing userId parameter",
		})
		log.Println("GET /users/:userId/certificates: Failed to get users certificates (missing userId parameter)")
		return
	}
	// Create new slice for the users certificates
	UserCertificates := []*Certificate{}
	// Check that the user is requesting certificates for a user they are allowed to (i.e. themself)
	if User.Id == RequestedID {
		// Loop through all certificates, add certificates that the user owns to the slice
		for _, Certificate := range Certificates {
			if Certificate.OwnerId == User.Id {
				UserCertificates = append(UserCertificates, Certificate)
			}
		}
	} else {
		c.JSON(401, map[string]string{
			"status": "failed",
			"error":  "Invalid User: You are requesting another users certificates",
		})
		log.Println("GET /users/:userId/certificates: Failed to get users certificates (requesting another users certificates)")
		return
	}
	// Respond with success message and certificates
	c.JSON(200, map[string]interface{}{
		"status":       "success",
		"certificates": UserCertificates,
	})
}

// Utility function to get the user from provided authorisation token or throw error
func GetUserOrError(c *gin.Context) *User {
	// Get Authorisation header value, or throw error
	Session := c.Request.Header.Get("Authorization")
	if Session == "" {
		c.JSON(401, map[string]string{
			"status": "failed",
			"error":  "No Authorization token",
		})
		return nil
	}
	// Check user exists, or throw error
	if UserID, Exists := Sessions[Session]; Exists {
		return Users[UserID]
	} else {
		c.JSON(401, map[string]string{
			"status": "failed",
			"error":  "Invalid Session",
		})
		return nil
	}
}
