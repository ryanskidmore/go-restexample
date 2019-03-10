package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Transfer struct {
	To     string `json:"to" binding:"required"`
	Status string `json:"status"`
}

func CreateTransferHandler(c *gin.Context) {
	User := GetUserOrError(c)
	if User == nil {
		return
	}
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
	CertID := c.Param("id")
	if CertID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Missing id parameter",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (missing id parameter)")
		return
	}
	TargetUserExists := false
	for _, User := range Users {
		if User.Email == TransferReq.To {
			TargetUserExists = true
		}
	}
	if !TargetUserExists {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Target User doesn't exist",
		})
		log.Println("POST /certificates/:id/transfers: Failed to create certificate transfer request (target user doesn't exist)")
		return
	}
	if Cert, Exists := Certificates[CertID]; Exists {
		if Cert.OwnerId == User.Id {
			TransferReq.Status = "valid"
			Certificates[CertID].Transfer = &TransferReq
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

func AcceptTransferHandler(c *gin.Context) {
	User := GetUserOrError(c)
	if User == nil {
		return
	}
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
	CertID := c.Param("id")
	if CertID == "" {
		c.JSON(400, map[string]string{
			"status": "failed",
			"error":  "Malformed Request: Missing id parameter",
		})
		log.Println("PUT /certificates/:id/transfers: Failed to update certificate transfer request (missing id parameter)")
		return
	}
	if Cert, Exists := Certificates[CertID]; Exists {
		if Cert.Transfer.To == User.Email {
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
			Certificates[CertID].Transfer = nil
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
