package permissions

import (
	"context"
	"net/http"
)

type contextKey string

const (
	SubjectContextKey contextKey = "hostmk_permissions_subject"
)

func WithSubject(
	ctx context.Context,
	subject *Subject,
) context.Context {

	if subject == nil {
		return ctx
	}

	return context.WithValue(
		ctx,
		SubjectContextKey,
		subject,
	)
}

func SubjectFromContext(
	ctx context.Context,
) (*Subject, bool) {

	if ctx == nil {
		return nil, false
	}

	subject, ok := ctx.Value(
		SubjectContextKey,
	).(*Subject)

	return subject, ok
}

type Middleware struct {
	checker *Checker
}

func NewMiddleware(
	checker *Checker,
) *Middleware {

	if checker == nil {
		checker = NewChecker()
	}

	return &Middleware{
		checker: checker,
	}
}

func (m *Middleware) RequirePermission(
	resource string,
	action string,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				subject, ok := SubjectFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrMissingUserContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if !m.checker.Can(
					subject,
					resource,
					action,
				) {

					http.Error(
						w,
						ErrPermissionDenied.Error(),
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

func (m *Middleware) RequireAnyPermission(
	resource string,
	actions ...string,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				subject, ok := SubjectFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrMissingUserContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if !m.checker.CanAny(
					subject,
					resource,
					actions...,
				) {

					http.Error(
						w,
						ErrPermissionDenied.Error(),
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

func (m *Middleware) RequireAllPermissions(
	resource string,
	actions ...string,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				subject, ok := SubjectFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrMissingUserContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if !m.checker.CanAll(
					subject,
					resource,
					actions...,
				) {

					http.Error(
						w,
						ErrPermissionDenied.Error(),
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

func (m *Middleware) RequireRole(
	roleName string,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				subject, ok := SubjectFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrMissingUserContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if !m.checker.HasRole(
					subject,
					roleName,
				) {

					http.Error(
						w,
						ErrPermissionDenied.Error(),
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

func (m *Middleware) RequireTenant(
	tenantID string,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				subject, ok := SubjectFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrMissingTenantContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if !m.checker.ValidateTenant(
					subject,
					tenantID,
				) {

					http.Error(
						w,
						ErrCrossTenantAccess.Error(),
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

func (m *Middleware) RequireSuperAdmin() func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				subject, ok := SubjectFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrMissingUserContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if !m.checker.IsSuperAdmin(
					subject,
				) {

					http.Error(
						w,
						ErrPermissionDenied.Error(),
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
