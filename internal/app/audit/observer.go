package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
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
}

func NewAuditSubject() *AuditSubject {
	return &AuditSubject{
		observers: make([]Observer, 0),
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

				fmt.Printf("Ошибка уведомления наблюдателя: %v\n", err)
			}
		}(observer)
	}
}

type FileAuditObserver struct {
	filePath string
	mutex    sync.Mutex
}

func NewFileAuditObserver(filePath string) *FileAuditObserver {
	return &FileAuditObserver{
		filePath: filePath,
	}
}

func (f *FileAuditObserver) Notify(ctx context.Context, event *AuditEvent) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла аудита: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ошибка сериализации события аудита: %w", err)
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("ошибка записи в файл аудита: %w", err)
	}

	return nil
}

type HTTPAuditObserver struct {
	url string
}

func NewHTTPAuditObserver(url string) *HTTPAuditObserver {
	return &HTTPAuditObserver{
		url: url,
	}
}

func (h *HTTPAuditObserver) Notify(ctx context.Context, event *AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ошибка сериализации события аудита: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("ошибка создания HTTP запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("сервер вернул ошибку: %d", resp.StatusCode)
	}

	return nil
}
