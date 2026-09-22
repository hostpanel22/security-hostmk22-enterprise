package fraud

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
	SaveSignal(ctx context.Context, signal *models.FraudSignal) (*models.FraudSignal, error)
	GetSignalByID(ctx context.Context, id string) (*models.FraudSignal, error)
	ListSignals(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudSignal, error)
	DeleteSignal(ctx context.Context, id string) error

	SaveViolation(ctx context.Context, violation *models.FraudViolation) (*models.FraudViolation, error)
	GetViolationByID(ctx context.Context, id string) (*models.FraudViolation, error)
	ListViolations(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudViolation, error)
	DeleteViolation(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	signals    map[string]*models.FraudSignal
	violations map[string]*models.FraudViolation
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		signals:    make(map[string]*models.FraudSignal),
		violations: make(map[string]*models.FraudViolation),
	}
}

func (r *InMemoryRepository) SaveSignal(ctx context.Context, signal *models.FraudSignal) (*models.FraudSignal, error) {
	if r == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if signal == nil {
		return nil, errors.New("fraud signal is required")
	}
	if err := signal.Validate(); err != nil {
		return nil, err
	}
	signal.Normalize()
	if strings.TrimSpace(signal.ID) == "" {
		signal.ID = "sig_" + time.Now().UTC().Format("20060102150405.000000000")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	clone := signal.Clone()
	r.signals[signal.ID] = &clone
	stored := r.signals[signal.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) GetSignalByID(ctx context.Context, id string) (*models.FraudSignal, error) {
	if r == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("fraud signal id is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	signal, ok := r.signals[id]
	if !ok || signal == nil {
		return nil, models.ErrNotFound
	}
	clone := signal.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) ListSignals(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudSignal, error) {
	if r == nil {
		return nil, errors.New("fraud repository is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*models.FraudSignal, 0, len(r.signals))
	for _, signal := range r.signals {
		if signal == nil {
			continue
		}
		if strings.TrimSpace(tenantID) != "" && signal.TenantID != tenantID {
			continue
		}
		clone := signal.Clone()
		items = append(items, &clone)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return paginateSignals(items, limit, offset), nil
}

func paginateSignals(items []*models.FraudSignal, limit, offset int) []*models.FraudSignal {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if offset >= len(items) {
		return []*models.FraudSignal{}
	}
	if limit == 0 || limit > len(items)-offset {
		limit = len(items) - offset
	}
	return items[offset : offset+limit]
}

func (r *InMemoryRepository) DeleteSignal(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("fraud signal id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.signals[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.signals, id)
	return nil
}

func (r *InMemoryRepository) SaveViolation(ctx context.Context, violation *models.FraudViolation) (*models.FraudViolation, error) {
	if r == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if violation == nil {
		return nil, errors.New("fraud violation is required")
	}
	if err := violation.Validate(); err != nil {
		return nil, err
	}
	violation.Normalize()
	if strings.TrimSpace(violation.ID) == "" {
		violation.ID = "vio_" + time.Now().UTC().Format("20060102150405.000000000")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	clone := violation.Clone()
	r.violations[violation.ID] = &clone
	stored := r.violations[violation.ID]
	result := stored.Clone()
	return &result, nil
}

func (r *InMemoryRepository) GetViolationByID(ctx context.Context, id string) (*models.FraudViolation, error) {
	if r == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("fraud violation id is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	violation, ok := r.violations[id]
	if !ok || violation == nil {
		return nil, models.ErrNotFound
	}
	clone := violation.Clone()
	return &clone, nil
}

func (r *InMemoryRepository) ListViolations(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudViolation, error) {
	if r == nil {
		return nil, errors.New("fraud repository is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*models.FraudViolation, 0, len(r.violations))
	for _, violation := range r.violations {
		if violation == nil {
			continue
		}
		if strings.TrimSpace(tenantID) != "" && violation.TenantID != tenantID {
			continue
		}
		clone := violation.Clone()
		items = append(items, &clone)
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
		return []*models.FraudViolation{}, nil
	}
	if limit == 0 || limit > len(items)-offset {
		limit = len(items) - offset
	}
	return items[offset : offset+limit], nil
}

func (r *InMemoryRepository) DeleteViolation(ctx context.Context, id string) error {
	if r == nil {
		return errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("fraud violation id is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.violations[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.violations, id)
	return nil
}
