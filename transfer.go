package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

// Transfer object
type Transfer struct {
	To     string `json:"to" binding:"required"`
	Status string `json:"status"`
}

// Handler to create a new transfer
func CreateTransferHandler(c *gin.Context) {
	// Get the user from provided authorisation token or throw error
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	// Bind body to transfer object
	var TransferReq Transfer
	err := c.BindJSON(&TransferReq)
	if err != nil {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (failed to bind post body to object: " + err.Error() + ")")
		return
	}
	// Get certificate ID from url parameters
	CertID := c.Param("id")
	if CertID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Missing id parameter",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (missing id parameter)")
		return
	}
	// Define boolean specifying if the target (to) user exists
	TargetUserExists := false
	// Loop through users and if user email matches the to field, set user exists boolean to true
	for _, User := range Users {
		if User.Email == TransferReq.To {
			TargetUserExists = true
			break // No need to loop through rest of users
		}
	}
	// If the target user doesn't exist
	if !TargetUserExists {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Target User doesn't exist",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (target user doesn't exist)")
		return
	}
	// Check certificate exists or throw error
	if Cert, Exists := Certificates[CertID]; Exists {
		// Check the user creating the transfer owns the certificate or throw error
		if Cert.OwnerId == User.Id {
			TransferReq.Status = "valid"
			// Assign transfer request to transfer field on certificate
			Certificates[CertID].Transfer = &TransferReq
			// Respond with success message
			c.JSON(200, map[string]string{
				"status": "success",
			})
		} else {
			c.JSON(401, map[string]string{
				"status": "failed",
				"error":  "Invalid User: You are not the owner of this certificate",
			})
			log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (user isn't owner of certificate)")
			return
		}
	} else {
		c.JSON(404, map[string]string{
			"status": "failed",
			"error":  "Certificate doesn't exist",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (certificate with ID doesn't exist)")
		return
	}
}

// Handler to accept a transfer
func AcceptTransferHandler(c *gin.Context) {
	// Get the user from provided authorisation token or throw error
	User := GetUserOrError(c)
	if User == nil {
		return
	}
	// Bind body to transfer object
	var TransferReq Transfer
	err := c.BindJSON(&TransferReq)
	if err != nil {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request",
		})
		log.Println("PUT /certificates/:id/transfers: Failed to update certificate transfer request (failed to bind post body to object: " + err.Error() + ")")
		return
	}
	// Get certificate id from url parameter
	CertID := c.Param("id")
	if CertID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Missing id parameter",
		})
		log.Println("PUT /certificates/:id/transfers: Failed to update certificate transfer request (missing id parameter)")
		return
	}
	// Check cert exists, or throw error
	if Cert, Exists := Certificates[CertID]; Exists {
		// Check the transfer target is the user accepting the transfer
		if Cert.Transfer.To == User.Email {
			// If status is accepted, complete the transfer, or if status is neither accepted or rejected then throw error
			if TransferReq.Status == "accepted" {
				Certificates[CertID].OwnerId = User.Id
			} else if TransferReq.Status != "accepted" && TransferReq.Status != "rejected" {
				c.JSON(400, map[string]string{
					"status": "failed",
					"error":  "Malformed Request: Invalid status (accepted or rejected)",
				})
				log.Println("PUT /certificates/:id/transfers: Failed to update certificate transfer request (invalid status)")
				return
			}
			// Set transfer field to nil (if status was rejected, this removes the transfer request - cancelling it)
			Certificates[CertID].Transfer = nil
			// Respond with success message
			c.JSON(200, map[string]string{
				"status": "success",
			})
		}
	} else {
		c.JSON(404, map[string]string{
			"status": "failed",
			"error":  "Certificate doesn't exist",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (certificate with ID doesn't exist)")
		return
	}
}
