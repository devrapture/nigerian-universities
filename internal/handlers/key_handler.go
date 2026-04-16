package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	apperrors "github.com/coolpythoncodes/nigerian-universities/internal/errors"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KeyHandlers struct {
	keyService service.KeyService
}

type GenerateKeyResponse struct {
	ID        uuid.UUID `json:"id"`
	Key       string    `json:"key"`
	IsActive  string    `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewKeyHandler(keyService service.KeyService) *KeyHandlers {
	return &KeyHandlers{
		keyService: keyService,
	}
}

func (h *KeyHandlers) CreateApiKey(c *gin.Context) {
	userID, _ := c.Get("userID")
	productKey, rawKey, err := h.keyService.HandleCreateKey(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to create api key")
		return
	}

	resp := GenerateKeyResponse{
		ID:        productKey.ID,
		Key:       rawKey,
		IsActive:  strconv.FormatBool(productKey.IsActive),
		CreatedAt: productKey.CreatedAt,
		UpdatedAt: productKey.UpdatedAt,
	}
	utils.SuccessResponse(c, http.StatusOK, "Store this key securely. It will not be shown again.", resp, nil)
}

func (h *KeyHandlers) GetAllKeys(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	keys, total, err := h.keyService.HandleGetAllKeys(c.Request.Context(), userID.(uuid.UUID), page, perPage)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to get all keys")
		return
	}

	meta := &utils.PaginationMeta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
		Pages:   int64(math.Ceil(float64(total) / float64(perPage))),
	}

	utils.SuccessResponse(c, http.StatusOK, "Fetched all keys", keys, meta)
}

func (h *KeyHandlers) RevokeKey(c *gin.Context) {
	userID, _ := c.Get("userID")
	keyID := c.Param("key_id")

	parsedKeyID, err := uuid.Parse(keyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "invalid key id")
		return
	}

	if err := h.keyService.HandleRevokeKey(c.Request.Context(), userID.(uuid.UUID), parsedKeyID); err != nil {
		switch err {
		case apperrors.ErrKeyNotFound:
			utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		case apperrors.ErrUnauthorized:
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to revoke key")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Key deactivated successfully", nil, nil)
}
