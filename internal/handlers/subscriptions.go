// Package handlers содержит http хендлеры
package handlers

import (
	"context"
	"net/http"
	"time"
	"errors"

	"subscription-service/internal/entities"
	"subscription-service/internal/logger"
	"subscription-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	repo repository.Repository
}

func New(repository repository.Repository) *Handler {
	return &Handler{repo: repository}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription record for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body entities.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} entities.Subscription
// @Failure 400 {object} entities.ErrorResponse "Invalid request"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Example subscription_create_request.json
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	log, ok := logger.GetLoggerFromCtx(c.Request.Context())
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Debug(c.Request.Context(), "create handler")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var req entities.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(c.Request.Context(), invalidJSONError, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.repo.CreateSubscription(ctx, req)
	if err != nil {
		log.Error(c.Request.Context(), "failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalError})
		return
	}

	log.Info(c.Request.Context(), "action=create resource=subscription status=success")
	c.JSON(http.StatusCreated, sub)
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Retrieve a specific subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 200 {object} entities.Subscription
// @Failure 404 {object} entities.ErrorResponse "Subscription not found"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	log, ok := logger.GetLoggerFromCtx(c.Request.Context())
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := getUUIDParam(c, "id")
	if err != nil {
		log.Warn(ctx, "get subscriptions", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.repo.GetSubscription(ctx, id)
	if err != nil {
		log.Error(ctx, "failed to get subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalError})
		return
	}

	if sub == nil {
		log.Warn(ctx, "subscription not found", zap.Any("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description Get a paginated list of subscriptions with optional filtering by user ID
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Filter by user ID" Format(uuid)
// @Param limit query int false "Number of records to return" default(10)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {array} entities.Subscription
// @Failure 400 {object} entities.ErrorResponse "Invalid request"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	log, ok := logger.GetLoggerFromCtx(c.Request.Context())
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*1e9)
	defer cancel()

	userID, err := getUUIDParam(c, "user_id")
	if err != nil {
		log.Error(ctx, "get subscriptions", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit, offset, err := getPaginationParams(c)
	if err != nil {
		log.Error(ctx, "Invalid beer id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subs, err := h.repo.ListSubscriptions(ctx, userID, limit, offset)
	if err != nil {
		log.Error(ctx, "failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if subs == nil {
		subs = []entities.Subscription{}
	}

	c.JSON(http.StatusOK, subs)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update an existing subscription (partial update supported)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Param request body entities.UpdateSubscriptionRequest true "Fields to update"
// @Success 200 {object} entities.Subscription
// @Failure 400 {object} entities.ErrorResponse "Invalid request"
// @Failure 404 {object} entities.ErrorResponse "Subscription not found"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	log, ok := logger.GetLoggerFromCtx(c.Request.Context())
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := getUUIDParam(c, "id")
	if err != nil {
		log.Error(ctx, "get subscriptions:", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req entities.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(c.Request.Context(), invalidJSONError, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.repo.UpdateSubscription(ctx, id, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			log.Error(c.Request.Context(), "Subs not found")
			c.JSON(http.StatusNotFound, gin.H{"error": notFoundError})
			return
		}
		log.Error(c.Request.Context(), "failed to update subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalError})
		return
	}

	if sub == nil {
		log.Warn(ctx, "subscription not found", zap.Any("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	log.Info(c.Request.Context(), "action=update resource=subscription status=success")
	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete a subscription by its ID
// @Tags subscriptions
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 204 "Subscription deleted successfully"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	log, ok := logger.GetLoggerFromCtx(c.Request.Context())
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := getUUIDParam(c, "id")
	if err != nil {
		log.Error(ctx, "get subscriptions", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.DeleteSubscription(ctx, id); err != nil {
		log.Error(ctx, "failed to delete subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalError})
		return
	}

	c.Status(http.StatusNoContent)
}


// CalculateCostReport godoc
// @Summary Calculate cost report
// @Description Calculate total cost of subscriptions for a given period with optional filtering
// @Tags reports
// @Accept json
// @Produce json
// @Param request body entities.CostReportRequest true "Report parameters"
// @Success 200 {object} entities.CostReport
// @Failure 400 {object} entities.ErrorResponse "Invalid request"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Router /reports/cost [post]
func (h *Handler) CalculateCostReport(c *gin.Context) {
	log, ok := logger.GetLoggerFromCtx(c.Request.Context())
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var req entities.CostReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, invalidJSONError, zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.repo.GetSubscriptionsCost(ctx, req)
	if err != nil {
		log.Error(ctx, "failed to calculate cost report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalError})
		return
	}

	c.JSON(http.StatusOK, report)
}
