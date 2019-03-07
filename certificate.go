package main

import (
	"github.com/gin-gonic/gin"
)

type Certificate struct {
	id        string
	title     string
	createdAt int64
	owner     *User
	year      int
	note      string
	transfer  *Transfer
}

func NewCertificateHandler(c *gin.Context) {

}

func UpdateCertificateHandler(c *gin.Context) {

}

func DeleteCertificateHandler(c *gin.Context) {

}
