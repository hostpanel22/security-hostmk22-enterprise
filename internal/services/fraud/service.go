package fraud

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/internal/models"
	repo "github.com/azezaz/HostMK-Security-Framework/internal/repositories/fraud"
)

type Service interface {
	SaveSignal(ctx context.Context, signal *models.FraudSignal) (*models.FraudSignal, error)
	GetSignalByID(ctx context.Context, id string) (*models.FraudSignal, error)
	ListSignals(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudSignal, error)
	DeleteSignal(ctx context.Context, id string) error

	SaveViolation(ctx context.Context, violation *models.FraudViolation) (*models.FraudViolation, error)
	GetViolationByID(ctx context.Context, id string) (*models.FraudViolation, error)
	ListViolations(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudViolation, error)
	DeleteViolation(ctx context.Context, id string) error
}

type service struct {
	repository repo.Repository
}

func NewService(repository repo.Repository) *service {
	return &service{repository: repository}
}

func (s *service) SaveSignal(ctx context.Context, signal *models.FraudSignal) (*models.FraudSignal, error) {
	if s == nil {
		return nil, errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if signal == nil {
		return nil, errors.New("fraud signal is required")
	}
	signal.Normalize()
	if strings.TrimSpace(signal.ID) == "" {
		signal.ID = "sig_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.SaveSignal(ctx, signal)
}

func (s *service) GetSignalByID(ctx context.Context, id string) (*models.FraudSignal, error) {
	if s == nil {
		return nil, errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("fraud signal id is required")
	}
	return s.repository.GetSignalByID(ctx, id)
}

func (s *service) ListSignals(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudSignal, error) {
	if s == nil {
		return nil, errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("fraud repository is nil")
	}
	return s.repository.ListSignals(ctx, tenantID, limit, offset)
}

func (s *service) DeleteSignal(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("fraud signal id is required")
	}
	return s.repository.DeleteSignal(ctx, id)
}

func (s *service) SaveViolation(ctx context.Context, violation *models.FraudViolation) (*models.FraudViolation, error) {
	if s == nil {
		return nil, errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if violation == nil {
		return nil, errors.New("fraud violation is required")
	}
	violation.Normalize()
	if strings.TrimSpace(violation.ID) == "" {
		violation.ID = "vio_" + time.Now().UTC().Format("20060102150405.000000000")
	}
	return s.repository.SaveViolation(ctx, violation)
}

func (s *service) GetViolationByID(ctx context.Context, id string) (*models.FraudViolation, error) {
	if s == nil {
		return nil, errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("fraud violation id is required")
	}
	return s.repository.GetViolationByID(ctx, id)
}

func (s *service) ListViolations(ctx context.Context, tenantID string, limit, offset int) ([]*models.FraudViolation, error) {
	if s == nil {
		return nil, errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return nil, errors.New("fraud repository is nil")
	}
	return s.repository.ListViolations(ctx, tenantID, limit, offset)
}

func (s *service) DeleteViolation(ctx context.Context, id string) error {
	if s == nil {
		return errors.New("fraud service is nil")
	}
	if s.repository == nil {
		return errors.New("fraud repository is nil")
	}
	if strings.TrimSpace(id) == "" {
		return errors.New("fraud violation id is required")
	}
	return s.repository.DeleteViolation(ctx, id)
}
