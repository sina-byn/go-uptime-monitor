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
	engine.Reload(true)

	app := fiber.New(fiber.Config{Views: engine})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		logs := db.ReadLogs()
		return c.Render("index", fiber.Map{"Logs": logs})
	})

	app.Post("/project", func(c *fiber.Ctx) error {
		project := db.Project{}

		if err := c.BodyParser(&project); err != nil {
			return err
		}

		db.CreateProject(project)

		return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "project added"})
	})

	c := cron.New()

	c.AddFunc("0 0 1 * *", func() {
		db.CleanupLogs()
	})

	c.AddFunc("* * * * *", func() {
		proejcts := db.ReadProjects()

		for _, project := range proejcts {
			go func() {
				resp, err := http.Get(project.Url)
				log := db.Log{Project: project.Name}

				if err != nil {
					log.Message = err.Error()
					log.Status = 500

					fmt.Println(err)
				} else {
					log.Message = "Success"
					log.Status = resp.StatusCode
				}

				db.CreateLog(log)
			}()
		}
	})

	db.InitDB()
	c.Start()
	app.Listen(":3000")
}
