Title: How to make a Simple Blog using Markdown and Go, Part 1
Date: 14 May 2024
Author: Mohd. Sameer
Intro: Welcome to the first blog post on my site! In this beginner-friendly series, I'll guide you through creating a blog using Go and Markdown. Assuming some prior programming experience in Git, CSS, and basic JavaScript, or any other language, the choice of Go and Markdown is justified by Go's simplicity and Markdown's easy syntax. Throughout this series, I'll guide you in setting up a basic Go web server, parsing Markdown files, templating for consistent styling, integrating CSS (I use Tailwind CSS), and HTMX, and adding features like search. By the end, you'll have a fully functional blog to showcase your writing. So, let's start!

### Prerequisites:

1. Go - [Download here](https://go.dev)
2. Text Editor - I use Vim BTW, you can use VS Code or any editor

### Let's get started!

Open your terminal in any folder and type

```
mkdir your_project_name
```
then go inside your folder
```
cd your_project_name
```
next type this in terminal
```
go mod init github.com/username/your_project_name
```
next
```
go get github.com/gofiber/fiber/v2
```
this will install go fiber framework, open it using your editor, for vim
```
vim .
```
for VS Code
```
code .
```
Great! Now let create our first file main.go

Press <kbd>%</kbd> in vim and write main.go in Vim, for VS Code type <kbd>Ctrl+N</kbd>, Now open it and type this
```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	
	app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
	})
	
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Server error:", err)
	}
}
```
Now if you open your browser and go to `localhost:3000` you will see
```
Hello, World!
```
Congrats! 🎉️🎉️🎉️ Go Server Done ✅️ Part 2 coming soon!
