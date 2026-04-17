package handlers

import (
	"math"
	"net/http"
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
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewKeyHandler(keyService service.KeyService) *KeyHandlers {
	return &KeyHandlers{
		keyService: keyService,
	}
}

// CreateApiKey creates a new api key for the user
// @Summary Create a new api key
// @Description Create a new api key for the user
// @Tags Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} schema.KeyCreateResponse
// @Failure 400 {object} schema.KeyBadRequestResponse
// @Failure 401 {object} schema.KeyUnauthorizedResponse
// @Failure 500 {object} schema.KeyInternalServerErrorResponse
// @Router /api-keys/generate [post]
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
		IsActive:  productKey.IsActive,
		CreatedAt: productKey.CreatedAt,
		UpdatedAt: productKey.UpdatedAt,
	}
	utils.SuccessResponse(c, http.StatusOK, "Store this key securely. It will not be shown again.", resp, nil)
}

// GetAllKeys returns a list of all keys for the user
// @Summary Get all keys
// @Description Get all keys for the user
// @Tags Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} schema.KeyListResponse
// @Failure 400 {object} schema.KeyBadRequestResponse
// @Failure 401 {object} schema.KeyUnauthorizedResponse
// @Failure 500 {object} schema.KeyInternalServerErrorResponse
// @Router /api-keys [get]
func (h *KeyHandlers) GetAllKeys(c *gin.Context) {
	userID, _ := c.Get("userID")
	queryDTO, err := parseListQuery(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	keys, total, err := h.keyService.HandleGetAllKeys(c.Request.Context(), userID.(uuid.UUID), queryDTO)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to get all keys")
		return
	}

	meta := &utils.PaginationMeta{
		Page:    queryDTO.Page,
		PerPage: queryDTO.Limit,
		Total:   total,
		Pages:   int64(math.Ceil(float64(total) / float64(queryDTO.Limit))),
	}

	utils.SuccessResponse(c, http.StatusOK, "Fetched all keys", keys, meta)
}

// RevokeKey deactivates a key
// @Summary Revoke a key
// @Description Deactivate a key
// @Tags Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key_id path string true "Key ID"
// @Success 200 {object} schema.KeySuccessResponse
// @Failure 400 {object} schema.KeyBadRequestResponse
// @Failure 401 {object} schema.KeyUnauthorizedResponse
// @Failure 404 {object} schema.KeyNotFoundResponse
// @Failure 500 {object} schema.KeyInternalServerErrorResponse
// @Router /api-keys/{key_id}/revoke [post]
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
