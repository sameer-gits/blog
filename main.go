package main

import (
    "html/template"
    "bytes"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html/v2"
    "github.com/yuin/goldmark"
)

func main() {
    engine := html.New("./views", ".html")

    app := fiber.New(fiber.Config{
        Views: engine,
    })

    app.Static("/", "./public")

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Get("posts/:post_name?", func(c *fiber.Ctx) error {
        post_name := c.Params("post_name")
        if post_name != "" {
            markdown, err := os.ReadFile("./all_blogs/" + post_name + ".md")
            if err != nil {
                log.Println("Error reading file:", err)
                return err
            }

            var post_content bytes.Buffer
            if err := goldmark.Convert(markdown, &post_content); err != nil {
                log.Println("Error converting markdown:", err)
                return err
            }

            return c.Render("post_templ", fiber.Map{
                "Title":   post_name,
                "Content": template.HTML(post_content.String()),
            })
        }

        return c.SendString("No post specified.")
    })

    if err := app.Listen(":3000"); err != nil {
        log.Fatal("Server error:", err)
    }
}
