package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"mcp-kodepos/internal/db"
	"mcp-kodepos/internal/mcp"
	"mcp-kodepos/internal/repo"
)

func main() {
	_ = godotenv.Load()

	dbase, err := db.New()
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer dbase.SQL.Close()

	r := gin.Default()

	// health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// MCP endpoint (JSON-RPC 2.0)
	h := mcp.NewHandler(repo.NewKodeposRepo(dbase.SQL))
	r.POST("/mcp", h.Post)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
