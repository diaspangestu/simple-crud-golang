package main

import (
	"belajar-golang-restful-api/controllers"
	"belajar-golang-restful-api/database"
	"belajar-golang-restful-api/helper"
	"belajar-golang-restful-api/middleware"
	"belajar-golang-restful-api/repository"
	"belajar-golang-restful-api/router"
	"belajar-golang-restful-api/services"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	db := database.NewDB()
	validate := validator.New()

	categoryRepository := repository.NewCategoryRepository()
	categoryService := services.NewCategoryService(categoryRepository, db, validate)
	categoryController := controllers.NewCategoryController(categoryService)

	route := router.NewRouter(categoryController)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: middleware.NewAuthMiddleware(route),
	}

	err := server.ListenAndServe()
	helper.PanicError(err)
}
