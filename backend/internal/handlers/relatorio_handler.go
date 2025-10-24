package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"relatorios-backend/internal/models"
	"relatorios-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type RelatorioHandler struct {
	relatorioService *services.RelatorioService
	pdfService       *services.PDFService
}

func NewRelatorioHandler(relatorioService *services.RelatorioService, pdfService *services.PDFService) *RelatorioHandler {
	return &RelatorioHandler{
		relatorioService: relatorioService,
		pdfService:       pdfService,
	}
}

func (h *RelatorioHandler) GetRelatorios(c *gin.Context) {
	// Retornar todos os relatórios (SEM AUTENTICAÇÃO)
	relatorios, err := h.relatorioService.GetAllRelatorios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"relatorios": relatorios})
}

func (h *RelatorioHandler) GetRelatorio(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do relatório é obrigatório"})
		return
	}

	relatorio, err := h.relatorioService.GetRelatorio(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, relatorio)
}

func (h *RelatorioHandler) CreateRelatorio(c *gin.Context) {
	var req models.CreateRelatorioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar dados: " + err.Error()})
		return
	}

	// Log para debug
	fmt.Printf("Criando relatório: %+v\n", req)

	relatorio, err := h.relatorioService.CreateRelatorio("system", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar relatório: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, relatorio)
}

func (h *RelatorioHandler) UpdateRelatorio(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do relatório é obrigatório"})
		return
	}

	var req models.UpdateRelatorioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	relatorio, err := h.relatorioService.UpdateRelatorio(id, "system", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, relatorio)
}

func (h *RelatorioHandler) DeleteRelatorio(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do relatório é obrigatório"})
		return
	}

	err := h.relatorioService.DeleteRelatorio(id, "system")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relatório deletado com sucesso"})
}

func (h *RelatorioHandler) GeneratePDF(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do relatório é obrigatório"})
		return
	}

	// Buscar relatório
	relatorio, err := h.relatorioService.GetRelatorio(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Removida verificação de propriedade - todos podem gerar PDF de qualquer relatório

	// Gerar PDF
	pdfBytes, err := h.pdfService.GeneratePDF(relatorio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar PDF: " + err.Error()})
		return
	}

	// Definir headers para download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=relatorio_"+id+".pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *RelatorioHandler) GenerateBatchPDF(c *gin.Context) {

	var req models.GenerateBatchPDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Por enquanto, retornar erro pois batch PDF é mais complexo
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Geração de PDF em lote ainda não implementada"})
}

// GenerateTestPDF - Rota pública para testar geração de PDF (REMOVER EM PRODUÇÃO)
func (h *RelatorioHandler) GenerateTestPDF(c *gin.Context) {
	// Criar um relatório de teste
	testRelatorio := &models.Relatorio{
		ID:           "test-123",
		Title:        "Relatório de Teste",
		TipoServico:  "MUTIRAO",
		Sub:          "CENTRO",
		Local:        "Praça da Liberdade",
		Endereco:     "Rua da Liberdade, 123",
		Descricao:    "Este é um relatório de teste para verificar se a geração de PDF está funcionando corretamente.",
		Data:         "2025-01-24",
		Fotos: []models.Foto{
			{URL: "https://via.placeholder.com/300x200?text=Foto+1", Descricao: "Foto de teste 1"},
			{URL: "https://via.placeholder.com/300x200?text=Foto+2", Descricao: "Foto de teste 2"},
		},
	}

	// Gerar PDF
	pdfBytes, err := h.pdfService.GeneratePDF(testRelatorio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar PDF: " + err.Error()})
		return
	}

	// Definir headers para download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=relatorio_teste.pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
