package main

import (
	"github.com/gin-gonic/gin"
)

type Transfer struct {
	To     *User
	Status string
}

func CreateTransferHandler(c *gin.Context) {

}

func AcceptTransferHandler(c *gin.Context) {

}
