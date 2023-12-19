package main

import (
	"context"
	"fmt"
	"gorm/models"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// type Interface interface {
// 	LogMode(LogLevel) Interface
// 	Info(context.Context, string, ...interface{})
// 	Warn(context.Context, string, ...interface{})
// 	Error(context.Context, string, ...interface{})
// 	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
// }

// We can implement only the method we need by embedding the logger.Interface on SqlLogger struct and overriding the Trace method.
// which is the only method we need to implement.
type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n====================\n", sql)
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}
}

func main() {
	initEnv()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dsn := fmt.Sprintf("%s?parseTime=True", os.Getenv("DATABASE_URL"))
	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: newLogger,
		// DryRun: true,
	})
	if err != nil {
		fmt.Printf("Error: %#v", err)
	}

	db.AutoMigrate(&models.Book{})
	fmt.Println("Database migrated")

	app := fiber.New()
	//GetBook
	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(models.GetBooks(db))
	})

	app.Get("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(models.GetBookByID(db, uint(id)))
	})

	app.Post("/books", func(c *fiber.Ctx) error {
		book := new(models.Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		err := models.CreateBook(db, book)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Created successfully",
		})
	})

	app.Put("/books/:id", func(c *fiber.Ctx) error {
		//Get id from url params
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book := new(models.Book)
		//Parse body into struct which BodyParser  need a pointer to struct(out parameter)
		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		//Update book
		book.ID = uint(id)
		err = models.UpdateBookByID(db, book)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		//Get latest book
		latestBook := models.GetBookByID(db, uint(id))
		return c.JSON(fiber.Map{
			"message": "Updated successfully",
			"status":  fiber.StatusOK,
			"data":    latestBook,
		})
	})

	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		//Get id from url params
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		//Delete book
		err = models.DeleteBookbyID(db, uint(id))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Deleted successfully",
		})
	})

	app.Listen(":8080")
}
