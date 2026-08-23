package notification

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Recipient, Subject, Body string
	CreatedAt                time.Time
}
type Sink interface {
	Send(context.Context, Message) error
}
type MemorySink struct {
	mu       sync.Mutex
	Messages []Message
	Failure  error
}

func (m *MemorySink) Send(ctx context.Context, msg Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Failure != nil {
		return m.Failure
	}
	m.Messages = append(m.Messages, msg)
	return nil
}
func (m *MemorySink) Snapshot() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Message, len(m.Messages))
	copy(out, m.Messages)
	return out
}

type Service struct{ Sink Sink }

func (s Service) Notify(ctx context.Context, to, subject string, parts ...string) error {
	if s.Sink == nil {
		return fmt.Errorf("notification sink missing")
	}
	body := ""
	for _, p := range parts {
		if body != "" {
			body += "\n"
		}
		body += p
	}
	return s.Sink.Send(ctx, Message{Recipient: to, Subject: subject, Body: body, CreatedAt: time.Now()})
}
