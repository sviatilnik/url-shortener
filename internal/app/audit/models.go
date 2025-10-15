package audit

import "time"

type AuditEvent struct {
	Timestamp int64  `json:"ts"`
	Action    string `json:"action"`
	UserID    string `json:"user_id"`
	URL       string `json:"url"`
}

func NewAuditEvent(action, userID, url string) *AuditEvent {
	return &AuditEvent{
		Timestamp: time.Now().Unix(),
		Action:    action,
		UserID:    userID,
		URL:       url,
	}
}

const (
	ActionShorten = "shorten"
	ActionFollow  = "follow"
)
