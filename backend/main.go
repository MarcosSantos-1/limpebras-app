package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"relatorios-backend/internal/config"
	"relatorios-backend/internal/database"
	"relatorios-backend/internal/handlers"
	"relatorios-backend/internal/middleware"
	"relatorios-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Carregar variáveis de ambiente
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Erro ao carregar variáveis de ambiente:", err)
	}

	// Conectar ao banco de dados
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.Close()

	// Executar migrações
	if err := database.Migrate(db); err != nil {
		log.Fatal("Erro ao executar migrações:", err)
	}

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configurar CORS
	corsOrigins := []string{"http://localhost:3000"}
	if corsEnv := os.Getenv("CORS_ORIGINS"); corsEnv != "" {
		// Dividir por vírgula e adicionar cada origem
		for _, origin := range strings.Split(corsEnv, ",") {
			corsOrigins = append(corsOrigins, strings.TrimSpace(origin))
		}
	}
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Inicializar serviços
	authService := services.NewAuthService(db)
	pdfService := services.NewPDFService()
	relatorioService := services.NewRelatorioService(db)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	relatorioHandler := handlers.NewRelatorioHandler(relatorioService, pdfService)

	// Todas as rotas são públicas (SEM AUTENTICAÇÃO)
	api := router.Group("/api")
	{
		// Relatórios
		api.GET("/relatorios", relatorioHandler.GetRelatorios)
		api.GET("/relatorios/:id", relatorioHandler.GetRelatorio)
		api.POST("/relatorios", relatorioHandler.CreateRelatorio)
		api.PUT("/relatorios/:id", relatorioHandler.UpdateRelatorio)
		api.DELETE("/relatorios/:id", relatorioHandler.DeleteRelatorio)
		
		// Geração de PDF
		api.POST("/relatorios/:id/pdf", relatorioHandler.GeneratePDF)
		api.POST("/relatorios/pdf/batch", relatorioHandler.GenerateBatchPDF)
		
		// Rota de teste para PDF
		api.POST("/test/pdf", relatorioHandler.GenerateTestPDF)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s", port)
	log.Fatal(router.Run(":" + port))
}
