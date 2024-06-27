package routes

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/shubhamgoyal1402/url-shortner/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	log.Printf("%s requested %s short url", c.IP(), url)

	rds0 := database.CreateClient(0)
	defer rds0.Close()

	value, err := rds0.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		log.Printf("%s requested short url - %s not found in database", c.IP(), url)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Error": "Short URL not found in the database",
		})
	} else if err != nil {
		log.Printf("%s request couldn't severed as service is unable to connect to the database", c.IP())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Can not connect to database",
		})
	}

	// Increment the counter
	rds1 := database.CreateClient(1)
	defer rds1.Close()
	_ = rds1.Incr(database.Ctx, "counter")

	// Redirect to original URL
	return c.Redirect(value, 301)
}
