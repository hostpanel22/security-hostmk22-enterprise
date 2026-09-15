package fraud

import (
	"context"
	"net/http"
)

type MiddlewareConfig struct {
	Detector *Detector

	Store Store

	BlockHighRisk bool
}

func FraudMiddleware(
	config MiddlewareConfig,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if config.Detector == nil {
					next.ServeHTTP(
						w,
						r,
					)
					return
				}

				signal := NewSignal(
					r.Header.Get("X-Tenant-ID"),
					r.Header.Get("X-User-ID"),
					"http_request",
					1,
					"middleware",
				)

				result, err := config.Detector.Detect(
					signal,
				)

				if err != nil {
					http.Error(
						w,
						err.Error(),
						http.StatusInternalServerError,
					)
					return
				}

				if config.Store != nil {

					_ = config.Store.SaveSignal(
						context.Background(),
						signal,
					)

					for _, violation := range result.Violations {

						_ = config.Store.SaveViolation(
							context.Background(),
							violation,
						)
					}
				}

				if config.BlockHighRisk &&
					result.Score.RequiresAction() {

					http.Error(
						w,
						"fraud risk detected",
						http.StatusForbidden,
					)

					return
				}

				next.ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}
