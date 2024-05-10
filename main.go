package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func renderPosts(directory string) ([]map[string]interface{}, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var posts []map[string]interface{}

	for _, file := range files {
		title := strings.TrimSuffix(file.Name(), ".md")
		capitalTitle := cases.Title(language.Und, cases.NoLower).String(title)
		markdown, err := os.ReadFile(directory + file.Name())
		if err != nil {
			log.Println("Error reading file:", err)
			continue
		}

		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
		)

		var postContent bytes.Buffer
		if err := md.Convert(markdown, &postContent); err != nil {
			log.Println("Error converting markdown:", err)
			continue
		}
		truncatedContent := postContent.String()
		post := map[string]interface{}{
			"Title":   capitalTitle,
            "Content": template.HTML(truncatedContent[:200]),
			"Slug":    strings.ReplaceAll(title, " ", "-"),
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

    directory := "./all_blogs/"

	app.Static("/", "./public")

	// Handler for the homepage to display all posts
	app.Get("/", func(c *fiber.Ctx) error {
		posts, err := renderPosts(directory)
		if err != nil {
			log.Println("Error reading blog directory:", err)
			return err
		}

		return c.Render("homepage", fiber.Map{
			"Posts": posts,
		})
	})

	// Handler for search
	app.Get("/search", func(c *fiber.Ctx) error {
		queryParam := c.Query("query")
		posts, err := renderPosts(directory)
		if err != nil {
			log.Println("Error reading blog directory:", err)
			return err
		}

		var filteredPosts []map[string]interface{}

		for _, post := range posts {
			title := post["Title"].(string)
			if queryParam == "" || strings.Contains(strings.ToLower(title), strings.ToLower(queryParam)) {
				filteredPosts = append(filteredPosts, post)
			}
		}

		return c.Render("post_homepage", fiber.Map{
			"Posts": filteredPosts,
		})
	})

	// Handler for individual posts
	app.Get("posts/:post_name?", func(c *fiber.Ctx) error {
		postName := c.Params("post_name")
		title := strings.ReplaceAll(postName, "-", " ")
		capitalTitle := cases.Title(language.Und, cases.NoLower).String(title)
		if postName != "" {
			markdown, err := os.ReadFile("./all_blogs/" + title + ".md")
			if err != nil {
				log.Println("Error reading file:", err)
				return err
			}
			md := goldmark.New(
				goldmark.WithExtensions(extension.GFM),
			)

			var postContent bytes.Buffer
			if err := md.Convert(markdown, &postContent); err != nil {
				log.Println("Error converting markdown:", err)
				return err
			}

			return c.Render("post_single", fiber.Map{
				"Title":   capitalTitle,
				"Content": template.HTML(postContent.String()),
			})
		}

		return c.SendString("No post specified.")
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Server error:", err)
	}
}
