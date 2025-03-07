package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/blockchain"
)

// Server represents the HTTP API server
type Server struct {
	blockchain *blockchain.Blockchain
	router    *gin.Engine
}

// NewServer creates a new HTTP API server
func NewServer(bc *blockchain.Blockchain) *Server {
	server := &Server{
		blockchain: bc,
		router:    gin.Default(),
	}
	server.setupRoutes()
	return server
}

// setupRoutes configures the API endpoints
func (s *Server) setupRoutes() {
	// Blockchain bilgisi
	s.router.GET("/info", s.getBlockchainInfo)
	s.router.GET("/blocks", s.getBlocks)
	s.router.GET("/blocks/:hash", s.getBlockByHash)
	
	// Validator işlemleri
	s.router.GET("/validators", s.getValidators)
	s.router.GET("/validators/current", s.getCurrentValidator)
	s.router.POST("/validators", s.addValidator)
	s.router.DELETE("/validators/:address", s.removeValidator)
	
	// İşlem gönderme
	s.router.POST("/transactions", s.submitTransaction)
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// getBlockchainInfo returns general information about the blockchain
func (s *Server) getBlockchainInfo(c *gin.Context) {
	info := gin.H{
		"blocks":     len(s.blockchain.Blocks),
		"validators": len(s.blockchain.Validators),
		"is_valid":   s.blockchain.IsValid(),
	}
	c.JSON(http.StatusOK, info)
}

// getBlocks returns all blocks in the chain
func (s *Server) getBlocks(c *gin.Context) {
	c.JSON(http.StatusOK, s.blockchain.Blocks)
}

// getBlockByHash returns a specific block by its hash
func (s *Server) getBlockByHash(c *gin.Context) {
	hash := c.Param("hash")
	for _, block := range s.blockchain.Blocks {
		if block.GetHashString() == hash {
			c.JSON(http.StatusOK, block)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
}

// getValidators returns all validators
func (s *Server) getValidators(c *gin.Context) {
	validators := make([]string, 0)
	for addr := range s.blockchain.Validators {
		validators = append(validators, addr)
	}
	c.JSON(http.StatusOK, validators)
}

// getCurrentValidator returns the current validator in rotation
func (s *Server) getCurrentValidator(c *gin.Context) {
	v, err := s.blockchain.GetCurrentValidator()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"address": v.Address})
}

// addValidator adds a new validator to the blockchain
func (s *Server) addValidator(c *gin.Context) {
	newValidator, err := validator.NewAuthority()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	s.blockchain.AddValidator(newValidator)
	c.JSON(http.StatusCreated, gin.H{
		"address": newValidator.Address,
		"message": "Validator added successfully",
	})
}

// removeValidator removes a validator from the blockchain
func (s *Server) removeValidator(c *gin.Context) {
	address := c.Param("address")
	s.blockchain.RemoveValidator(address)
	c.JSON(http.StatusOK, gin.H{
		"message": "Validator removed successfully",
	})
}

// TransactionRequest represents a new transaction request
type TransactionRequest struct {
	Data      string `json:"data" binding:"required"`
	Validator string `json:"validator" binding:"required"`
}

// submitTransaction submits a new transaction to be added to the blockchain
func (s *Server) submitTransaction(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validator'ı bul
	v, exists := s.blockchain.Validators[req.Validator]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validator not found"})
		return
	}

	// İşlemi blockchain'e ekle
	err := s.blockchain.AddBlock(req.Data, v)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction added successfully",
		"block":   s.blockchain.Blocks[len(s.blockchain.Blocks)-1],
	})
} 