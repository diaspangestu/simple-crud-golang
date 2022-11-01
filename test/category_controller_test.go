package test

import (
	"belajar-golang-restful-api/controllers"
	"belajar-golang-restful-api/helper"
	"belajar-golang-restful-api/middleware"
	"belajar-golang-restful-api/models/entity"
	"belajar-golang-restful-api/repository"
	"belajar-golang-restful-api/router"
	"belajar-golang-restful-api/services"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:veronica123@tcp(localhost:3306)/belajar_golang_restful_api_test")
	helper.PanicError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter() http.Handler {
	db := setupTestDB()
	validate := validator.New()

	categoryRepository := repository.NewCategoryRepository()
	categoryService := services.NewCategoryService(categoryRepository, db, validate)
	categoryController := controllers.NewCategoryController(categoryService)

	route := router.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(route)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	route := setupRouter()
	requestBody := strings.NewReader(`{"name": "TESTING"}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "TESTING", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)
	route := setupRouter()
	requestBody := strings.NewReader(`{"name": ""}`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Create(context.Background(), tx, entity.Category{
		Name: "Gadget",
	})
	tx.Commit()

	route := setupRouter()
	requestBody := strings.NewReader(`{"name": "TESTING"}`)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "TESTING", responseBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Create(context.Background(), tx, entity.Category{
		Name: "Gadget",
	})
	tx.Commit()

	route := setupRouter()
	requestBody := strings.NewReader(`{"name": ""}`)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Create(context.Background(), tx, entity.Category{
		Name: "Gadget",
	})
	tx.Commit()

	route := setupRouter()
	requestBody := strings.NewReader(`{"name": "Gadget"}`)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	route := setupRouter()
	requestBody := strings.NewReader(`{"name": ""}`)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/categories/404", requestBody)
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Create(context.Background(), tx, entity.Category{
		Name: "Gadget",
	})
	tx.Commit()

	route := setupRouter()

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	route := setupRouter()

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/categories/404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestGetAllCategorySuccess(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category1 := categoryRepository.Create(context.Background(), tx, entity.Category{
		Name: "Gadget",
	})
	category2 := categoryRepository.Create(context.Background(), tx, entity.Category{
		Name: "Makanan",
	})
	tx.Commit()

	route := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/categories", nil)
	request.Header.Add("X-API-Key", "SECRET")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	var categories = responseBody["data"].([]interface{})

	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])

	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()
	truncateCategory(db)

	route := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/categories", nil)
	request.Header.Add("X-API-Key", "SALAH")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "Unauthorized", responseBody["status"])
}
