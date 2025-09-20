package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
var rdb *redis.Client
var ctx = context.Background()

func main() {
	// Инициализация базы данных
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация Redis
	rdb = initRedis()

	// Проверяем подключение к Redis
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v (continuing without cache)", err)
		rdb = nil
	} else {
		log.Println("Redis connected successfully")
	}

	// Автоматическая миграция
	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	// Существующие роуты
	r.GET("/ping", PingPong)
	r.GET("/redis", RedisStatus)

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

// initRedis инициализирует подключение к Redis
func initRedis() *redis.Client {
	// Пробуем подключиться к redis (Docker), если не получается - к localhost
	host := "redis"
	if !isHostReachable(host + ":6379") {
		host = "localhost"
	}

	return redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // без пароля
		DB:       0,  // база данных по умолчанию
	})
} // isHostReachable проверяет доступность хоста
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

	// Очищаем кэш книг после создания новой
	if rdb != nil {
		rdb.Del(ctx, "books:all")
	}

	c.JSON(http.StatusCreated, book)
} // GetAllBooks возвращает все книги с кэшированием
func GetAllBooks(c *gin.Context) {
	cacheKey := "books:all"

	// Пытаемся получить из кэша если Redis доступен
	if rdb != nil {
		cached, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var books []Book
			if json.Unmarshal([]byte(cached), &books) == nil {
				c.Header("X-Cache", "HIT")
				c.JSON(http.StatusOK, books)
				return
			}
		}
	}

	// Получаем из базы данных
	var books []Book
	if result := db.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	// Сохраняем в кэш на 5 минут если Redis доступен
	if rdb != nil {
		if booksJSON, err := json.Marshal(books); err == nil {
			rdb.Set(ctx, cacheKey, booksJSON, 5*time.Minute)
		}
	}

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, books)
}

func PingPong(c *gin.Context) {
	c.String(200, "pong")
}

func RedisStatus(c *gin.Context) {
	if rdb == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unavailable",
			"message": "Redis not connected",
		})
		return
	}

	// Проверяем подключение
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "connected",
		"message": "Redis is working",
	})
}
