package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/robfig/cron/v3"
	"github.com/sina-byn/go-uptime-monitoring/internal/db"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		logs := db.ReadLogs()
		return c.Render("index", fiber.Map{"Logs": logs})
	})

	c := cron.New()

	c.AddFunc("*/1 * * * *", func() {
		db.CleanupLogs()
	})

	c.AddFunc("*/1 * * * *", func() {
		resp, err := http.Get("https://aia.tools")
		log := db.Log{Project: "aia"}

		if err != nil {
			log.Message = "Failed"
			log.Status = 500

			fmt.Println(err)
		} else {
			log.Message = resp.Status
			log.Status = resp.StatusCode
		}

		db.CreateLog(log)
	})

	db.InitDB()
	c.Start()
	app.Listen(":3000")
}
