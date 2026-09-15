package fraud_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/fraud"
)

func TestMiddlewareWithoutDetector(t *testing.T) {
	handler := fraud.FraudMiddleware(
		fraud.MiddlewareConfig{},
	)(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}

func TestMiddlewareAllowRequest(t *testing.T) {
	detector := fraud.NewDetector()

	store := fraud.NewMemoryStore()

	handler := fraud.FraudMiddleware(
		fraud.MiddlewareConfig{
			Detector: detector,
			Store:    store,
		},
	)(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		"X-Tenant-ID",
		"tenant-1",
	)

	req.Header.Set(
		"X-User-ID",
		"user-1",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}

	if store.SignalCount() != 1 {
		t.Fatal("signal should be stored")
	}
}

func TestMiddlewareBlockHighRisk(t *testing.T) {
	detector := fraud.NewDetector()

	rule := fraud.NewRule(
		"high-risk",
		"block request",
		"http_request",
		1,
		90,
	)

	_ = detector.AddRule(
		rule,
	)

	handler := fraud.FraudMiddleware(
		fraud.MiddlewareConfig{
			Detector:      detector,
			BlockHighRisk: true,
		},
	)(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		"X-Tenant-ID",
		"tenant-1",
	)

	req.Header.Set(
		"X-User-ID",
		"user-1",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected %d got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}
