package main

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

var Certificates map[string]*Certificate

type Certificate struct {
	Id        string    `json:"id"`
	Title     string    `json:"title" binding:"required"`
	CreatedAt int64     `json:"createdAt"`
	OwnerId   string    `json:"ownerId"`
	Year      int       `json:"year" binding:"required"`
	Note      string    `json:"note"`
	Transfer  *Transfer `json:"transfer"`
}

func NewCertificateHandler(c *gin.Context) {
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	var Cert Certificate
	err := c.BindJSON(&Cert)
	if err != nil {
		c.AbortWithStatus(400)
		log.Println("POST /certificates: Failed to add new certificate (failed to bind post body to object: " + err.Error() + ")")
		return
	}
	RandUUID, err := uuid.NewV4()
	if err != nil {
		c.AbortWithStatus(500)
		log.Println("POST /certificates: Failed to add new certificate (failed to generate UUIDv4: " + err.Error() + ")")
		return
	}
	Cert.Id = RandUUID.String()
	Cert.CreatedAt = time.Now().Unix()
	Cert.OwnerId = User.Id
	Certificates[Cert.Id] = &Cert
	log.Println("POST /certificates: Successfully created new certificate (" + Cert.Id + ": " + Cert.Title + ")")
	c.JSON(200, map[string]string{
		"status": "success",
		"id":     Cert.Id,
	})
}

func UpdateCertificateHandler(c *gin.Context) {
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	var UpdatedCert Certificate
	err := c.BindJSON(&UpdatedCert)
	if err != nil {
		c.AbortWithStatus(400)
		log.Println("PUT /certificates/:id: Failed to update certificate (failed to bind post body to object: " + err.Error() + ")")
		return
	}
	CertID := c.Param("id")
	if CertID == "" {
		c.AbortWithStatus(400)
		log.Println("PUT /certificates/:id: Failed to update certificate (missing id Parameter)")
		return
	}
	if Cert, Exists := Certificates[CertID]; Exists {
		if Cert.OwnerId == User.Id {
			Cert.Title = UpdatedCert.Title
			Cert.Year = UpdatedCert.Year
			Cert.Note = UpdatedCert.Note
			Certificates[CertID] = Cert
			log.Println("PUT /certificates/:id: Successfully updated certificate (" + Cert.Id + ": " + Cert.Title + ")")
			c.JSON(200, map[string]string{
				"status": "success",
				"id":     Cert.Id,
			})
		} else {
			c.AbortWithStatus(401)
			log.Println("PUT /certificates/:id: Failed to update certificate (user isn't owner of certificate)")
			return
		}
	} else {
		c.AbortWithStatus(404)
		log.Println("PUT /certificates/:id: Failed to update certificate (certificate with ID doesn't exist)")
		return
	}
}

func DeleteCertificateHandler(c *gin.Context) {
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	CertID := c.Param("id")
	if CertID == "" {
		c.AbortWithStatus(400)
		log.Println("DELETE /certificates/:id: Failed to delete certificate (missing id Parameter)")
		return
	}
	if Cert, Exists := Certificates[CertID]; Exists {
		if Cert.OwnerId == User.Id {
			delete(Certificates, CertID)
			log.Println("DELETE /certificates/:id: Successfully deleted certificate (" + Cert.Id + ": " + Cert.Title + ")")
			c.JSON(200, map[string]string{
				"status": "success",
				"id":     Cert.Id,
			})
		} else {
			c.AbortWithStatus(401)
			log.Println("DELETE /certificates/:id: Failed to delete certificate (user isn't owner of certificate)")
			return
		}
	} else {
		c.AbortWithStatus(404)
		log.Println("DELETE /certificates/:id: Failed to delete certificate (certificate with id doesn't exist)")
		return
	}
}
