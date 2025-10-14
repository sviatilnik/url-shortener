package audit

import (
	"context"

	"go.uber.org/zap"
)

type AuditService struct {
	subject Subject
	log     *zap.SugaredLogger
}

func NewAuditService(log *zap.SugaredLogger) *AuditService {
	return &AuditService{
		subject: NewAuditSubject(log),
		log:     log,
	}
}

func (s *AuditService) AddFileObserver(filePath string) error {
	if filePath == "" {
		return nil
	}

	observer := NewFileAuditObserver(filePath, s.log)
	s.subject.Attach(observer)
	return nil
}

func (s *AuditService) AddHTTPObserver(url string) error {
	if url == "" {
		return nil
	}

	observer := NewHTTPAuditObserver(url, s.log)
	s.subject.Attach(observer)
	return nil
}

func (s *AuditService) LogEvent(ctx context.Context, action, userID, url string) {
	event := NewAuditEvent(action, userID, url)
	s.subject.NotifyObservers(ctx, event)
}

func (s *AuditService) LogShortenEvent(ctx context.Context, userID, url string) {
	s.LogEvent(ctx, ActionShorten, userID, url)
}

func (s *AuditService) LogFollowEvent(ctx context.Context, userID, url string) {
	s.LogEvent(ctx, ActionFollow, userID, url)
}
