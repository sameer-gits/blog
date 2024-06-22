package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	fiberHTML "github.com/gofiber/template/html/v2"
	"github.com/russross/blackfriday/v2"
)

const (
	primaryColor = "cyan"
	directory    = "./all_blogs/"
	dateLayout   = "02 Jan 2006"
)

type Work struct {
	ID          int     `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	Video       *string `json:"video"`
	Link        *string `json:"link"`
}

var dbpool *pgxpool.Pool

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Printf("DATABASE_URL environment variable is not set")
	}

	dbpool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
	}
}

func readMarkdownFile(filePath string) (map[string]string, string, error) {
	markdown, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	msg, err := mail.ReadMessage(bytes.NewReader(markdown))
	if err != nil {
		return nil, "", err
	}

	metadata := map[string]string{
		"Date":   msg.Header.Get("Date"),
		"Author": msg.Header.Get("Author"),
		"Title":  msg.Header.Get("Title"),
		"Intro":  msg.Header.Get("Intro"),
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(msg.Body)
	if err != nil {
		return nil, "", err
	}

	return metadata, buf.String(), nil
}

func searchPosts(directory string) ([]map[string]interface{}, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var posts []map[string]interface{}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		title := strings.TrimSuffix(file.Name(), ".md")
		filePath := filepath.Join(directory, file.Name())

		metadata, _, err := readMarkdownFile(filePath)
		if err != nil {
			log.Println("Error reading file:", err)
			continue
		}

		post := map[string]interface{}{
			"Metadata":     metadata,
			"Slug":         strings.ReplaceAll(title, " ", "-"),
			"PrimaryColor": primaryColor,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func renderPosts(directory string) ([]map[string]interface{}, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var posts []map[string]interface{}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		title := strings.TrimSuffix(file.Name(), ".md")
		filePath := filepath.Join(directory, file.Name())

		metadata, markdownContent, err := readMarkdownFile(filePath)
		if err != nil {
			log.Println("Error reading file:", err)
			continue
		}

		body := blackfriday.Run([]byte(markdownContent))

		post := map[string]interface{}{
			"Content":      template.HTML(body),
			"Slug":         strings.ReplaceAll(title, " ", "-"),
			"PrimaryColor": primaryColor,
			"Metadata":     metadata,
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		dateI, errI := time.Parse(dateLayout, posts[i]["Metadata"].(map[string]string)["Date"])
		dateJ, errJ := time.Parse(dateLayout, posts[j]["Metadata"].(map[string]string)["Date"])

		if errI != nil || errJ != nil {
			log.Println("Error parsing dates:", errI, errJ)
			return false
		}

		return dateI.After(dateJ)
	})

	return posts, nil
}

func main() {
	defer dbpool.Close()

	engine := fiberHTML.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	// Handler for the homepage to display all posts
	app.Get("/", func(c *fiber.Ctx) error {
		posts, err := renderPosts(directory)
		if err != nil {
			log.Println("Error reading blog directory:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading blog directory")
		}

		return c.Render("homepage", fiber.Map{
			"Posts":        posts,
			"PrimaryColor": primaryColor,
		})
	})

	app.Get("/portfolio", func(c *fiber.Ctx) error {
		rows, err := dbpool.Query(context.Background(), "SELECT id, title, description, image, link, video FROM work")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch work from database")
		}
		defer rows.Close()

		var works []Work
		for rows.Next() {
			var work Work
			err := rows.Scan(&work.ID, &work.Title, &work.Description, &work.Image, &work.Link, &work.Video)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to scan work data: %v", err.Error()))
			}
			works = append(works, work)
		}

		if rows.Err() != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to iterate over work rows")
		}
		return c.Render("portfolio", fiber.Map{
			"Works": works,
		})
	})

	// Handler for search
	app.Get("/search", func(c *fiber.Ctx) error {
		queryParam := strings.ToLower(c.Query("query"))
		posts, err := searchPosts(directory)
		if err != nil {
			log.Println("Error reading blog directory:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading blog directory")
		}

		var filteredPosts []map[string]interface{}

		for _, post := range posts {
			metadata := post["Metadata"].(map[string]string)
			title := metadata["Title"]

			if queryParam == "" {

			} else if strings.Contains(strings.ToLower(title), strings.ToLower(queryParam)) {
				// If the title contains the query string, include the post
				filteredPosts = append(filteredPosts, post)
			}
		}

		return c.Render("post_results", fiber.Map{
			"FilteredPosts": filteredPosts,
			"PrimaryColor":  primaryColor,
		})
	})

	// Handler for individual posts
	app.Get("/posts/:post_name?", func(c *fiber.Ctx) error {
		postName := c.Params("post_name")
		if postName != "" {
			title := strings.ReplaceAll(postName, "-", " ")
			filePath := filepath.Join(directory, title+".md")

			metadata, markdownContent, err := readMarkdownFile(filePath)
			if err != nil {
				log.Println("Error reading file:", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
			}

			body := blackfriday.Run([]byte(markdownContent))

			return c.Render("post_single", fiber.Map{
				"Content":      template.HTML(body),
				"Metadata":     metadata,
				"PrimaryColor": primaryColor,
			})
		}

		return c.SendStatus(fiber.StatusBadRequest)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Fatal("Server error:", err)
	}
}
