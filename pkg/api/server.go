package api

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/blockchain"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
)

// Server represents the HTTP API server
type Server struct {
	blockchain *blockchain.Blockchain
	router    *gin.Engine
}

// NewServer creates a new HTTP API server
func NewServer(bc *blockchain.Blockchain) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	
	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Upgrade, Connection")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Upgrade")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Request logging
		fmt.Printf("[API] %s %s - Headers: %v\n", c.Request.Method, c.Request.URL.Path, c.Request.Header)
		
		c.Next()

		// Response logging
		fmt.Printf("[API] Response Status: %d\n", c.Writer.Status())
	})

	server := &Server{
		blockchain: bc,
		router:    router,
	}
	server.setupRoutes()
	return server
}

// setupRoutes configures the API endpoints
func (s *Server) setupRoutes() {
	// WebSocket endpoint
	s.router.GET("/ws", func(c *gin.Context) {
		s.blockchain.WebSocketServer.HandleWebSocket(c.Writer, c.Request)
	})

	// Blockchain bilgisi
	s.router.GET("/info", s.getBlockchainInfo)
	s.router.GET("/blocks", s.getBlocks)
	s.router.GET("/blocks/:hash", s.getBlockByHash)
	s.router.GET("/transactions", s.getTransactions)
	
	// Validator işlemleri
	s.router.GET("/validators", s.getValidators)
	s.router.GET("/validators/current", s.getCurrentValidator)
	s.router.POST("/validators", s.addValidator)
	s.router.DELETE("/validators/:address", s.removeValidator)
	
	// İşlem gönderme
	s.router.POST("/transactions", s.submitTransaction)

	// Akıllı kontrat endpoint'leri
	contracts := s.router.Group("/contracts")
	{
		contracts.POST("", s.deployContract)
		contracts.GET("", s.listContracts)
		contracts.GET("/:address", s.getContract)
		contracts.POST("/:address/execute", s.executeContract)
		contracts.POST("/:address/disable", s.disableContract)
		contracts.POST("/:address/enable", s.enableContract)
	}
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// getBlockchainInfo returns general information about the blockchain
func (s *Server) getBlockchainInfo(c *gin.Context) {
	info := gin.H{
		"blocks":              len(s.blockchain.Blocks),
		"validators":          len(s.blockchain.Validators),
		"is_valid":           s.blockchain.IsValid(),
		"current_block":      len(s.blockchain.Blocks) - 1,
		"active_validators":  s.blockchain.GetActiveValidatorCount(),
		"validator_count":    len(s.blockchain.Validators),
		"pending_transactions": 0, // TODO: Add mempool integration
	}

	fmt.Printf("[API] Sending blockchain info: %+v\n", info)
	c.JSON(http.StatusOK, info)
}

// getBlocks returns all blocks in the chain
func (s *Server) getBlocks(c *gin.Context) {
	blocks := make([]gin.H, 0)
	for _, block := range s.blockchain.Blocks {
		blocks = append(blocks, gin.H{
			"height": block.Header.Height,
			"hash": block.GetHashString(),
			"timestamp": block.Header.Timestamp.Unix() * 1000, // Convert to milliseconds
			"transactions": len(block.Transactions),
			"validator": block.Header.ValidatorAddress,
		})
	}
	c.JSON(http.StatusOK, blocks)
}

// getTransactions returns all transactions from recent blocks
func (s *Server) getTransactions(c *gin.Context) {
	transactions := make([]gin.H, 0)
	
	// Get transactions from the last 10 blocks
	startBlock := 0
	if len(s.blockchain.Blocks) > 10 {
		startBlock = len(s.blockchain.Blocks) - 10
	}
	
	for _, block := range s.blockchain.Blocks[startBlock:] {
		for _, tx := range block.Transactions {
			transactions = append(transactions, gin.H{
				"hash": hex.EncodeToString(tx.Hash),
				"from": tx.From,
				"to": tx.To,
				"value": tx.Value.String(),
				"type": "transfer", // Default to transfer type for now
				"timestamp": block.Header.Timestamp.Unix() * 1000,
				"status": "success", // Default to success for now
			})
		}
	}
	
	c.JSON(http.StatusOK, transactions)
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
	// Create a new private key for the validator
	newValidator, err := validator.NewAuthority(nil)
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

	// Find the validator
	v, exists := s.blockchain.Validators[req.Validator]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validator not found"})
		return
	}

	// Create a new block
	prevBlock := s.blockchain.GetLatestBlock()
	newBlock, err := blockchain.NewBlock(
		prevBlock.Header.Height + 1,
		prevBlock.GetHash(),
		prevBlock.Header.StateRoot,
		1000000, // Default gas limit
		v,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add the block to the chain
	if err := s.blockchain.AddBlock(newBlock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction added successfully",
		"block":   newBlock,
	})
}

// deployContract handles contract deployment
func (s *Server) deployContract(c *gin.Context) {
	var req struct {
		Code    string `json:"code" binding:"required"`
		Owner   string `json:"owner" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Version string `json:"version" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kontrat kodunu decode et
	code, err := hex.DecodeString(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract code"})
		return
	}

	// Owner adresini parse et
	owner := common.HexToAddress(req.Owner)

	// Kontratı deploy et
	contract, err := s.blockchain.DeployContract(code, owner, req.Name, req.Version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

// listContracts handles contract listing
func (s *Server) listContracts(c *gin.Context) {
	contracts := s.blockchain.ListContracts()
	c.JSON(http.StatusOK, contracts)
}

// getContract handles contract retrieval
func (s *Server) getContract(c *gin.Context) {
	address := common.HexToAddress(c.Param("address"))
	contract, err := s.blockchain.GetContract(address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// executeContract handles contract execution
func (s *Server) executeContract(c *gin.Context) {
	address := common.HexToAddress(c.Param("address"))

	var req struct {
		Input string `json:"input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input verisini decode et
	input, err := hex.DecodeString(req.Input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	// Kontratı çalıştır
	result, err := s.blockchain.ExecuteContract(address, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": hex.EncodeToString(result),
	})
}

// disableContract handles contract disabling
func (s *Server) disableContract(c *gin.Context) {
	address := common.HexToAddress(c.Param("address"))

	var req struct {
		Owner string `json:"owner" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner := common.HexToAddress(req.Owner)
	if err := s.blockchain.ContractManager.DisableContract(address, owner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// enableContract handles contract enabling
func (s *Server) enableContract(c *gin.Context) {
	address := common.HexToAddress(c.Param("address"))

	var req struct {
		Owner string `json:"owner" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner := common.HexToAddress(req.Owner)
	if err := s.blockchain.ContractManager.EnableContract(address, owner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
} 