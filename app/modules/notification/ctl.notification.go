package notification

import (
	"app/app/helper"
	notificationdto "app/app/modules/notification/dto"
	"app/app/response"
	"app/internal/logger"
	"log"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		Service: svc,
	}
}

func (ctl *Controller) Connect(ctx *gin.Context) {
	user, _ := helper.GetUserByToken(ctx)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Register client with the hub
	client := &Client{
		Hub:    ctl.Service.Hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: user.Data.ID,
	}

	client.Hub.Register <- client

	// Start goroutines
	go client.WritePump()
	go client.ReadPump()
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := notificationdto.ListNotificationRequest{
		Page:    1,
		Size:    10,
		OrderBy: "asc",
		SortBy:  "created_at",
	}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

// MarkAsRead marks a notification as read
func (ctl *Controller) MarkAsRead(ctx *gin.Context) {
	ID := notificationdto.GetByIDNotification{}
	if err := ctx.ShouldBindUri(&ID); err != nil {
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	user, _ := helper.GetUserByToken(ctx)

	err := ctl.Service.MarkAsRead(ctx, ID.ID, user.Data.ID)
	if err != nil {
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, nil)
}

// MarkAllAsRead marks all notifications as read for current user
func (ctl *Controller) MarkAllAsRead(ctx *gin.Context) {
	user, _ := helper.GetUserByToken(ctx)

	err := ctl.Service.MarkAllAsRead(ctx, user.Data.ID)
	if err != nil {
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, nil)
}

// GetUnreadCount gets count of unread notifications
func (ctl *Controller) GetUnreadCount(ctx *gin.Context) {
	user, _ := helper.GetUserByToken(ctx)

	count, err := ctl.Service.GetUnreadCount(ctx, user.Data.ID)
	if err != nil {
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, count)
}
