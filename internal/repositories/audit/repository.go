package audit

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
)

type Repository interface {
	Create(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error)
	GetByID(ctx context.Context, id string) (*models.AuditEvent, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]*models.AuditEvent, error)
	Update(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	events map[string]*models.AuditEvent
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{events: make(map[string]*models.AuditEvent)}
}

func (r *InMemoryRepository) Create(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error) {
	if r == nil {
		return nil, errors.New("audit repository is nil")
	}
	if event == nil {
		return nil, errors.New("audit event is required")
	}
	if err := event.Validate(); err != nil {
		return nil, err
	}
	event.Normalize()
	if strings.TrimSpace(event.ID) == "" {
		event.ID = "evt_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	if event.UpdatedAt.IsZero() {
		event.UpdatedAt = event.CreatedAt
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	copyEvent := event.Clone()
	r.events[event.ID] = &copyEvent
	stored := r.events[event.ID]
	clone := stored.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*models.AuditEvent, error) {
	if r == nil {
		return nil, errors.New("audit repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("audit event id is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	event, ok := r.events[id]
	if !ok || event == nil {
		return nil, models.ErrNotFound
	}
	clone := event.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]*models.AuditEvent, error) {
	if r == nil {
		return nil, errors.New("audit repository is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*models.AuditEvent, 0, len(r.events))
	for _, event := range r.events {
		if event == nil {
			continue
		}
		if strings.TrimSpace(tenantID) != "" && event.TenantID != tenantID {
			continue
		}
		clone := event.Clone()
		items = append(items, &clone)
	}
	if len(items) == 0 {
		return []*models.AuditEvent{}, nil
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if offset >= len(items) {
		return []*models.AuditEvent{}, nil
	}
	if limit == 0 || limit > len(items)-offset {
		limit = len(items) - offset
	}
	return items[offset : offset+limit], nil
}

func (r *InMemoryRepository) Update(ctx context.Context, event *models.AuditEvent) (*models.AuditEvent, error) {
	if r == nil {
		return nil, errors.New("audit repository is nil")
	}
	if event == nil {
		return nil, errors.New("audit event is required")
	}
	if err := event.Validate(); err != nil {
		return nil, err
	}
	event.Normalize()

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.events[event.ID]; !ok {
		return nil, models.ErrNotFound
	}
	clone := event.Clone()
	r.events[event.ID] = &clone
	stored := r.events[event.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("audit repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("audit event id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.events[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.events, id)
	return nil
}
