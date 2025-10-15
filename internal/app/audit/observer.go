package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"
)

type Observer interface {
	Notify(ctx context.Context, event *AuditEvent) error
}

type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	NotifyObservers(ctx context.Context, event *AuditEvent)
}

type AuditSubject struct {
	observers []Observer
	mutex     sync.RWMutex
	log       *zap.SugaredLogger
}

func NewAuditSubject(log *zap.SugaredLogger) *AuditSubject {
	return &AuditSubject{
		observers: make([]Observer, 0),
		log:       log,
	}
}

func (s *AuditSubject) Attach(observer Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *AuditSubject) Detach(observer Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *AuditSubject) NotifyObservers(ctx context.Context, event *AuditEvent) {
	s.mutex.RLock()
	observers := make([]Observer, len(s.observers))
	copy(observers, s.observers)
	s.mutex.RUnlock()

	for _, observer := range observers {
		go func(obs Observer) {
			if err := obs.Notify(ctx, event); err != nil {
				s.log.Errorw("Ошибка уведомления наблюдателя", "error", err)
			}
		}(observer)
	}
}

type FileAuditObserver struct {
	filePath string
	mutex    sync.Mutex
	log      *zap.SugaredLogger
}

func NewFileAuditObserver(filePath string, log *zap.SugaredLogger) *FileAuditObserver {
	return &FileAuditObserver{
		filePath: filePath,
		log:      log,
	}
}

func (f *FileAuditObserver) Notify(ctx context.Context, event *AuditEvent) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f.log.Errorw("Ошибка открытия файла аудита", "error", err)
		return err
	}
	defer file.Close()

	data, err := json.Marshal(event)
	if err != nil {
		f.log.Errorw("Ошибка сериализации события аудита", "error", err)
		return err
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		f.log.Errorw("Ошибка записи в файл аудита", "error", err)
		return err
	}

	return nil
}

type HTTPAuditObserver struct {
	url string
	log *zap.SugaredLogger
}

func NewHTTPAuditObserver(url string, log *zap.SugaredLogger) *HTTPAuditObserver {
	return &HTTPAuditObserver{
		url: url,
		log: log,
	}
}

func (h *HTTPAuditObserver) Notify(ctx context.Context, event *AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		h.log.Errorw("Ошибка сериализации события аудита", "error", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.url, bytes.NewBuffer(data))
	if err != nil {
		h.log.Errorw("Ошибка создания HTTP запроса", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Errorw("Ошибка отправки HTTP запроса", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.log.Errorw("Сервер вернул ошибку", "status_code", resp.StatusCode)
		return fmt.Errorf("сервер вернул ошибку: %d", resp.StatusCode)
	}

	return nil
}
