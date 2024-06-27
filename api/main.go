package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/shubhamgoyal1402/url-shortner/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("form.html")
	})

	app.Get("/:url", routes.ResolveURL)
	app.Post("/shorten", handleForm)
	app.Post("/api/v1", routes.ShortenURL)
}

func handleForm(c *fiber.Ctx) error {
	url := c.FormValue("url")
	customShort := c.FormValue("custom")
	expiry := c.FormValue("expiry")

	// Create a request body to send to the ShortenURL route
	body := routes.Request{
		URL:         url,
		CustomShort: customShort,
	}

	if expiry != "" {
		exp, err := strconv.Atoi(expiry)
		if err == nil {
			body.Expiry = time.Duration(exp)
		}
	}

	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Unable to process request",
		})
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/api/v1", bytes.NewBuffer(reqBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Unable to create request",
		})
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Unable to process request",
		})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return c.JSON(result)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()
	app.Static("/", "./")
	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}
