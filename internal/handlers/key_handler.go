package handlers

import (
	"log"

	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/gin-gonic/gin"
)

type KeyHandlers struct {
	keyService service.KeyService
}

func NewKeyHandler(keyService service.KeyService) *KeyHandlers {
	return &KeyHandlers{
		keyService: keyService,
	}
}

func (h *KeyHandlers) CreateApiKey(c *gin.Context) {
	userID, _ := c.Get("userID")
	log.Println("userID", userID)
	// h.keyService.HandleCreateKey(c.Request.Context(), userID.(uuid.UUID))
}
