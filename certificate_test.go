package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestNewCertificateHandler(t *testing.T) {
	// Initialise api server
	go main()
	// Wait for server to start
	time.Sleep(5 * time.Second)
	// Create new certificate
	TestCert := Certificate{
		Title: "Example Work",
		Year:  1950,
		Note:  "Example Note",
	}
	// Convert new cert into JSON
	TestCertJSON, err := json.Marshal(TestCert)
	if err != nil {
		t.Error("Failed to create new certificate, error during JSON marshal: " + err.Error())
		return
	}
	// Create new request to POST/create certificate
	Request, err := http.NewRequest("POST", "http://localhost:8080/certificates", bytes.NewBuffer(TestCertJSON))
	if err != nil {
		t.Error("Failed to create new certificate, error while creating new request: " + err.Error())
		return
	}
	// Set appropiate headers
	Request.Header.Set("Authorization", "2b65701217e2c546bbf69a982c50014d503b977be64b7687616858ef8b6ed411")
	Request.Header.Set("Content-Type", "application/json")
	Client := &http.Client{}
	// Perform request
	Response, err := Client.Do(Request)
	if err != nil {
		t.Error("Failed to create new certificate, error while performing request: " + err.Error())
		return
	}
	defer Response.Body.Close()
	// Read response
	RespBody, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		t.Error("Failed to create new certificate, error while reading response: " + err.Error())
		return
	}
	// Convert response into map
	var RespMap map[string]interface{}
	err = json.Unmarshal(RespBody, &RespMap)
	if err != nil {
		t.Error("Failed to create new certificate, error while decoding response: " + err.Error())
		return
	}
	// Check if certificate creation was successful
	CertStatus := RespMap["status"].(string)
	if CertStatus == "failed" {
		t.Error("Failed to create new certificate: " + RespMap["error"].(string))
		return
	}
	// Check cert exists
	CertID := RespMap["id"].(string)
	if Cert, Exists := Certificates[CertID]; Exists {
		if Cert.Title == TestCert.Title && Cert.Year == TestCert.Year && Cert.Note == TestCert.Note {
			t.Log("Data matched")
			return
		} else {
			t.Error("Failed to create new certificate: data mismatch between submitted and saved")
			return
		}
	} else {
		t.Error("Failed to create new certificate: certificate doesn't exist in map")
		return
	}
}

func TestUpdateCertificateHandler(t *testing.T) {

}

func TestDeleteCertificateHandler(t *testing.T) {

}
