package handler

import (
	"net/http"

	"delayed-notifier/internal/model"
	"delayed-notifier/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.NotificationService
}

func (h *Handler) CreateNotification(c *gin.Context) {
	var req model.Notification

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.Service.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
