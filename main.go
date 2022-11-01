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
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
)

func main() {
	PORT := os.Getenv("PORT")
	db := database.NewDB()
	validate := validator.New()

	categoryRepository := repository.NewCategoryRepository()
	categoryService := services.NewCategoryService(categoryRepository, db, validate)
	categoryController := controllers.NewCategoryController(categoryService)

	route := router.NewRouter(categoryController)

	server := http.Server{
		Addr:    "/" + PORT,
		Handler: middleware.NewAuthMiddleware(route),
	}

	err := server.ListenAndServe()
	helper.PanicError(err)
}
