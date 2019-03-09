package main

import (
	"github.com/gin-gonic/gin"
)

var Sessions map[string]string

func main() {
	InitDummyData()
	router := gin.Default()
	router.POST("/certificates", NewCertificateHandler)
	router.PUT("/certificates/:id", UpdateCertificateHandler)
	router.DELETE("/certificates/:id", DeleteCertificateHandler)
	router.GET("/users/:userId/certificates", GetCertificatesHandler)
	router.POST("/certificates/:id/transfers", CreateTransferHandler)
	router.PUT("/certificates/:id/transfers", AcceptTransferHandler)
	router.Run()
}

func InitDummyData() {
	Certificates = make(map[string]*Certificate)
	Users = make(map[string]*User)
	Sessions = make(map[string]string)
	Users["db75530b-34d4-47df-bb08-82f075b6045b"] = &User{
		Id:    "db75530b-34d4-47df-bb08-82f075b6045b",
		Email: "bob@examplemail.com",
		Name:  "Bob",
	}
	Users["11150398-01a1-4335-afde-4cbaa5d490f9"] = &User{
		Id:    "11150398-01a1-4335-afde-4cbaa5d490f9",
		Email: "alice@examplemail.com",
		Name:  "Alice",
	}
	Sessions["2b65701217e2c546bbf69a982c50014d503b977be64b7687616858ef8b6ed411"] = "db75530b-34d4-47df-bb08-82f075b6045b"
	Sessions["9609301053e5cb2b2e43258e163120a70f420adfacda30c5d8d1e069ae02f028"] = "11150398-01a1-4335-afde-4cbaa5d490f9"
}
