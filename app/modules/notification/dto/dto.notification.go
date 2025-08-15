package notificationdto

type NotificationRequest struct {
	Action     string `json:"action"`
	UserID     string `json:"user_id"`
	Message    string `json:"message"`
	DocumentID string `json:"document_id"`
}

type NotificationRead struct {
	Type string `json:"type"`
	Data struct {
		NotificationID string `json:"notificationId"`
		IsRead         bool   `json:"isRead"`
		UserID         string `json:"userId"`
		Count          int    `json:"count"`
	} `json:"data"`
}

type GetByIDNotification struct {
	ID string `uri:"id" binding:"required"`
}

type ListNotificationRequest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"searchBy"`
	SortBy   string `form:"sortBy"`
	OrderBy  string `form:"orderBy"`
	UserID   string `form:"userId"`
	IsRead   *bool  `form:"isRead"`
}
