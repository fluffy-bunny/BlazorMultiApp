package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/xid"
)

var (
	indexRoot string
	indexApp1 string
	indexApp2 string
)

func noCacheMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasSuffix(c.Request().URL.Path, "index.html") ||
			c.Request().URL.Path == "/" ||
			c.Request().URL.Path == "/app1" ||
			c.Request().URL.Path == "/app2" {
			c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Response().Header().Set("Pragma", "no-cache")
			c.Response().Header().Set("Expires", "0")
		}
		return next(c)
	}
}

func serveIndexFromMemory(content string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(http.StatusOK, content)
	}
}

func createAppHandler(staticMiddleware echo.MiddlewareFunc, indexContent string, rootPath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if path == rootPath {
			return c.HTML(http.StatusOK, indexContent)
		}

		// This is likely a file request, try to serve it statically
		err := staticMiddleware(func(c echo.Context) error {
			return nil
		})(c)
		if err == nil {
			return nil // File was found and served
		}
		return c.NoContent(http.StatusNotFound)
	}

}

func main() {
	port := flag.String("port", "9044", "Port to listen on")
	flag.Parse()

	guid := xid.New().String()

	// Load index files into memory
	indexRootContent, err := os.ReadFile("static/index_template.html")
	if err != nil {
		panic(err)
	}
	indexApp1Content, err := os.ReadFile("static/app1/wwwroot/index_template.html")
	if err != nil {
		panic(err)
	}
	indexApp2Content, err := os.ReadFile("static/app2/wwwroot/index_template.html")
	if err != nil {
		panic(err)
	}

	// Replace {version} with guid in all index files
	indexRoot = strings.ReplaceAll(string(indexRootContent), "{version}", guid)
	indexApp1 = strings.ReplaceAll(string(indexApp1Content), "{version}", guid)
	indexApp2 = strings.ReplaceAll(string(indexApp2Content), "{version}", guid)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(noCacheMiddleware)

	// Serve root index.html from memory
	e.GET("/", serveIndexFromMemory(indexRoot))

	// Serve static files and handle routing for app1
	app1Static := middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static/app1/wwwroot",
		HTML5:  true,
		Browse: false,
	})

	e.GET("/app1/*", createAppHandler(app1Static, indexApp1, "/app1/"))

	// Serve static files and handle routing for app2
	app2Static := middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static/app2/wwwroot",
		HTML5:  true,
		Browse: false,
	})

	e.GET("/app2/*", createAppHandler(app2Static, indexApp2, "/app2/"))

	// Serve other static files from the root static folder
	e.Static("/", "static")

	fmt.Printf("Server started on port %s\n", *port)
	e.Logger.Fatal(e.Start(":" + *port))
}
