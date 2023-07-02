package main

import (
	"net/http"
	"os"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email		string  	`json:"email"`
	Status		bool  		`json:"status"`
}

type Todo struct {
	gorm.Model
	// UserId	int			`json:"user_id"`
	// User		User        `json:"user"`
	Name  		string		`json:"name"`
	Completed	bool		`json:"completed"`
}

type DeleteTodo struct {
	TodoId		int			`json:"todo_id"`
}	

type UpdateTodo struct {
	TodoId		int			`json:"todo_id"`
	Name		string		`json:"name"`
	Completed	bool		`json:"completed"`
}	

func main() {

	// dsn := "root:@tcp(127.0.0.1:3306)/gohome?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:@tcp(host.docker.internal:3306)/gohome?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&Todo{}, &User{})

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/create_todo", func(c echo.Context) error {
		newTodo := new(Todo)
		if err := c.Bind(&newTodo); err != nil {
			log.Fatalln(err)
			return c.JSON(http.StatusBadRequest, err)
		}

		createError := db.Create(&newTodo).Error
		if createError != nil{
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.DELETE("/delete_todo", func(c echo.Context) error {
		deleteTodo := new(DeleteTodo)
		if err := c.Bind(&deleteTodo); err != nil {
			log.Fatalln(err)
			return c.JSON(http.StatusBadRequest, err)
		}

		log.Print(deleteTodo)

		deleteError := db.Delete(&Todo{}, deleteTodo.TodoId).Error
		if deleteError != nil{
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})


	e.GET("/todo", func(c echo.Context) error {
		var todos []Todo
		if res := db.Find(&todos); res.Error != nil {
			data := map[string]interface{}{
				"message": res.Error.Error(),
			}
			return c.JSON(http.StatusOK, data)
		}

		response := map[string]interface{}{
			"data": todos,
		}

		return c.JSON(http.StatusOK, response)
	})

	e.PUT("/update_todo", func(c echo.Context) error {
		updateTodo := new(UpdateTodo)
		if err := c.Bind(&updateTodo); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		var todo Todo
		searchError := db.First(&todo, updateTodo.TodoId).Error
		if searchError != nil{
			return c.JSON(http.StatusNotFound, err)
		}

		todo.Name = updateTodo.Name 
		todo.Completed = updateTodo.Completed
		updateError := db.Save(todo).Error
		if updateError != nil{
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Simple implementation of an integer minimum
// Adapted from: https://gobyexample.com/testing-and-benchmarking
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}