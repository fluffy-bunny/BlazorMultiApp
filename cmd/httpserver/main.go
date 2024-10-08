package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func serveIndexHandler(rootPath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		indexPath := rootPath + "/index.html"
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "Index file not found")
		}
		return c.File(indexPath)
	}
}
func serveModifiedIndexHandler(modifiedContent string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(http.StatusOK, modifiedContent)
	}
}
func createAppHandler(staticMiddleware echo.MiddlewareFunc, rootPath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := staticMiddleware(func(c echo.Context) error {
			return nil
		})(c); err == nil {
			return nil
		}
		return serveIndexHandler(rootPath)(c)
	}
}

func main() {
	port := flag.String("port", "9044", "Port to listen on")
	flag.Parse()
	/*
		guid := xid.New().String()

		indexApp1, err := os.ReadFile("static/app1/wwwroot/index_template.html")
		if err != nil {
			panic(err)
		}
		indexApp2, err := os.ReadFile("static/app2/wwwroot/index_template.html")
		if err != nil {
			panic(err)
		}

		// Convert the content to a string
		contentStr := string(indexApp1)

		// Replace all instances of {version} with "guid"
		modifiedApp1 := strings.ReplaceAll(contentStr, "{version}", guid)
		// Convert the content to a string
		contentStr = string(indexApp2)

		// Replace all instances of {version} with "guid"
		modifiedApp2 := strings.ReplaceAll(contentStr, "{version}", guid)
	*/
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(noCacheMiddleware)

	// Serve root index.html
	e.GET("/", serveIndexHandler("static"))

	// Serve static files and handle routing for app1
	app1Static := middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static/app1/wwwroot",
		HTML5:  true,
		Browse: false,
	})
	e.GET("/app1", serveIndexHandler("static/app1/wwwroot"))
	e.GET("/app1/*", createAppHandler(app1Static, "static/app1/wwwroot"))

	// Serve static files and handle routing for app2
	app2Static := middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static/app2/wwwroot",
		HTML5:  true,
		Browse: false,
	})
	e.GET("/app2", serveIndexHandler("static/app2/wwwroot"))
	e.GET("/app2/*", createAppHandler(app2Static, "static/app2/wwwroot"))

	// Serve other static files from the root static folder
	e.Static("/", "static")

	fmt.Printf("Server started on port %s\n", *port)
	e.Logger.Fatal(e.Start(":" + *port))
}
