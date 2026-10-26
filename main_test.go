package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)


func SetUpRouter() *gin.Engine{
	router := gin.Default()
	return router
}


func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome to the Tech Company listing API with Golang. Updated version 2"}`

	r := SetUpRouter()
	r.GET("/", HomepageHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)


	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewCompanyHandler(t *testing.T) {
	r := SetUpRouter()
	r.POST("/company", NewCompanyHandler)
	companyId := xid.New().String()

	company := Company{
		ID: companyId,
		Name: "Demo Company",
		CEO: "Demo CEO",
		Revenue: "35 million",
	}

	jsonValue, _ := json.Marshal(company)
	req, _ := http.NewRequest("POST", "/company", bytes.NewBuffer(jsonValue))
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetCompaniesHandler(t *testing.T) {
	r := SetUpRouter()
	r.GET("/companies", GetCompaniesHandler)
	req, _ := http.NewRequest("GET", "/companies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var companies []Company
	json.Unmarshal(w.Body.Bytes(), &companies)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, companies)
}



func TestUpdateCompanyHandler(t *testing.T) {
	r := SetUpRouter()
	r.PUT("/company/:id", UpdateCompanyHandler)

	company := Company{
		ID: "3",
		Name: "Microsoft",
		CEO: "Satya Nadella",
		Revenue: "245 billion",
	}

	jsonValue, _ := json.Marshal(company)
	req, _ := http.NewRequest("PUT", "/company/3", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	for _, c := range companies {
		if c.ID == "3" {
			assert.Equal(t, "245 billion", c.Revenue)
		}
		if c.ID == "2" {
			assert.Equal(t, "Netflix", c.Name)
		}
	}
}

func TestDeleteCompanyHandler(t *testing.T) {
	r := SetUpRouter()
	r.DELETE("/company/:id", DeleteCompanyHandler)
	req, _ := http.NewRequest("DELETE", "/company/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	for _, c := range companies {
		assert.NotEqual(t, "1", c.ID)
	}
}
