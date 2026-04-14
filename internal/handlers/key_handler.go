package handlers

import (
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KeyHandlers interface{}

type keyHandlers struct {
	keyService service.KeyService
}

func NewKeyHandler(keyService service.KeyService) KeyHandlers {
	return &keyHandlers{
		keyService: keyService,
	}
}

func (h *keyHandlers) CreateApiKey(c *gin.Context) {
	userID, _ := c.Get("userID")
	h.keyService.HandleCreateKey(c.Request.Context(), userID.(uuid.UUID))
}
