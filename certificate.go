package main

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

// Map of certificates
// Certificate ID -> *Certificate
var Certificates map[string]*Certificate

// Certificate object
type Certificate struct {
	Id        string    `json:"id"`
	Title     string    `json:"title" binding:"required"`
	CreatedAt int64     `json:"createdAt"`
	OwnerId   string    `json:"ownerId"`
	Year      int       `json:"year" binding:"required"`
	Note      string    `json:"note"`
	Transfer  *Transfer `json:"transfer"`
}

// Handler for creating a new certifcate
func NewCertificateHandler(c *gin.Context) {
	// Get the user from provided authorisation token or throw error
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	// Bind body to certificate object
	var Cert Certificate
	err := c.BindJSON(&Cert)
	if err != nil {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request",
		})
		log.Println("POST /certificates: Failed to add new certificate (failed to bind post body to object: " + err.Error() + ")")
		return
	}
	// Generate new random UUID for certificate
	RandUUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(500, map[string]string{
			"status": "failed",
			"error":  "Internal Server Error",
		})
		log.Println("POST /certificates: Failed to add new certificate (failed to generate UUIDv4: " + err.Error() + ")")
		return
	}
	// Assign non-provided/generated fields
	Cert.Id = RandUUID.String()
	Cert.CreatedAt = time.Now().Unix()
	Cert.OwnerId = User.Id
	// Add certificate to map
	Certificates[Cert.Id] = &Cert
	log.Println("POST /certificates: Successfully created new certificate (" + Cert.Id + ": " + Cert.Title + ")")
	// Respond with success message
	c.JSON(200, map[string]string{
		"status": "success",
		"id":     Cert.Id,
	})
}

// Handler for updating a certificate
func UpdateCertificateHandler(c *gin.Context) {
	// Get the user from provided authorisation token or throw error
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	// Bind body to certificate object
	var UpdatedCert Certificate
	err := c.BindJSON(&UpdatedCert)
	if err != nil {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request",
		})
		log.Println("PUT /certificates/:id: Failed to update certificate (failed to bind post body to object: " + err.Error() + ")")
		return
	}
	// Get Certificate ID from URL parameters
	CertID := c.Param("id")
	if CertID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Missing id parameter",
		})
		log.Println("PUT /certificates/:id: Failed to update certificate (missing id parameter)")
		return
	}
	// Check if cert exists, and throw error if not
	if Cert, Exists := Certificates[CertID]; Exists {
		// Check the user owns the certificate they're trying to update, or throw an error
		if Cert.OwnerId == User.Id {
			// Update the certificate fields that the user should be allowed to update
			Cert.Title = UpdatedCert.Title
			Cert.Year = UpdatedCert.Year
			Cert.Note = UpdatedCert.Note
			// Put certificate back into map
			Certificates[CertID] = Cert
			log.Println("PUT /certificates/:id: Successfully updated certificate (" + Cert.Id + ": " + Cert.Title + ")")
			// Respond with success message
			c.JSON(200, map[string]string{
				"status": "success",
				"id":     Cert.Id,
			})
		} else {
			c.JSON(401, map[string]string{
				"status": "failed",
				"error":  "Invalid User: You are not the owner of this certificate",
			})
			log.Println("PUT /certificates/:id: Failed to update certificate (user isn't owner of certificate)")
			return
		}
	} else {
		c.JSON(404, map[string]string{
			"status": "failed",
			"error":  "Certificate doesn't exist",
		})
		log.Println("PUT /certificates/:id: Failed to update certificate (certificate with ID doesn't exist)")
		return
	}
}

// Handler for deleting a certificate
func DeleteCertificateHandler(c *gin.Context) {
	// Get the user from provided authorisation token or throw error
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	// Get Certificate ID from URL parameters
	CertID := c.Param("id")
	if CertID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Missing id parameter",
		})
		log.Println("DELETE /certificates/:id: Failed to delete certificate (missing id parameter)")
		return
	}
	// Check if cert exists, and throw error if not
	if Cert, Exists := Certificates[CertID]; Exists {
		// Check the user owns the certificate they're trying to update, or throw an error
		if Cert.OwnerId == User.Id {
			// Delete certificate from map
			delete(Certificates, CertID)
			log.Println("DELETE /certificates/:id: Successfully deleted certificate (" + Cert.Id + ": " + Cert.Title + ")")
			// Respond with success message
			c.JSON(200, map[string]string{
				"status": "success",
				"id":     Cert.Id,
			})
		} else {
			c.JSON(401, map[string]string{
				"status": "failed",
				"error":  "Invalid User: You are not the owner of this certificate",
			})
			log.Println("DELETE /certificates/:id: Failed to delete certificate (user isn't owner of certificate)")
			return
		}
	} else {
		c.JSON(404, map[string]string{
			"status": "failed",
			"error":  "Certificate doesn't exist",
		})
		log.Println("DELETE /certificates/:id: Failed to delete certificate (certificate with id doesn't exist)")
		return
	}
}
