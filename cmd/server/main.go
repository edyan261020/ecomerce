package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. Servir archivos estáticos (CSS, JS, imágenes)
	r.Static("/static", "./static")

	// 2. Cargar plantillas HTML
	r.LoadHTMLGlob("templates/*")

	// 3. Página de bienvenida
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"title": "Bienvenido a mi e-comerce",
		})
	})

	// 4. Página de productos
	r.GET("/products", func(c *gin.Context) {
		c.HTML(http.StatusOK, "products.html", gin.H{
			"title": "Productos",
		})
	})

	// 6. Página de login
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Iniciar Sesión",
		})
	})

	// 7. Página "Acerca de"
	r.GET("/sobre", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sobre.html", gin.H{
			"title": "Acerca de e-comerce",
		})
	})
	// 7. Página "Acerca de"
	r.GET("/adminproducts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminproducts.html", gin.H{
			"title": "Acerca de e-comerce",
		})
	})
	r.GET("/adminusers", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminusers.html", gin.H{
			"title": "Acerca de e-comerce",
		})
	})
	r.GET("/adminorders", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminorders.html", gin.H{
			"title": "Acerca de e-comerce",
		})
	})
	r.GET("/adminreports", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminreports.html", gin.H{
			"title": "Acerca de e-comerce",
		})
	})

	// 8. Iniciar el servidor
	r.Run(":8080") // http://localhost:8080
}
