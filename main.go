package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Book модель для книги
type Book struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Author    string    `json:"author" gorm:"not null"`
	ISBN      string    `json:"isbn" gorm:"unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *gorm.DB

func main() {
	// Инициализация базы данных
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автоматическая миграция
	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	// Существующие роуты
	r.GET("/ping", PingPong)
	r.GET("/newping", NewPong)

	// Новые роуты для работы с книгами
	r.POST("/books", CreateBook)
	r.GET("/books", GetAllBooks)

	if err := r.Run(); err != nil {
		panic(err)
	}
}

// initDB инициализирует подключение к PostgreSQL
func initDB() (*gorm.DB, error) {
	// Пробуем подключиться к postgres (Docker), если не получается - к localhost
	host := "postgres"
	if !isHostReachable(host + ":5432") {
		host = "localhost"
	}

	dsn := "host=" + host + " port=5432 user=postgres password=postgres dbname=cicd_learn sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// isHostReachable проверяет доступность хоста
func isHostReachable(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// CreateBook создает новую книгу
func CreateBook(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Create(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetAllBooks возвращает все книги
func GetAllBooks(c *gin.Context) {
	var books []Book

	if result := db.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func PingPong(c *gin.Context) {
	c.String(200, "pong")
}

func NewPong(c *gin.Context) {
	c.String(200, "new pong")
}
