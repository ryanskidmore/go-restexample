package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/certificates", NewCertificateHandler)
	router.PUT("/certificates/:id", UpdateCertificateHandler)
	router.DELETE("/certificates/:id", DeleteCertificateHandler)
	router.GET("/users/:userId/certificates", GetCertificatesHandler)
	router.POST("/certificates/:id/transfers", CreateTransferHandler)
	router.PUT("/certificates/:id/transfers", AcceptTransferHandler)
	router.Run()
}
