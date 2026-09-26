package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	iamclient "order/internal/client/iam"
)

type contextKey string

const (
	userContextKey    contextKey = "auth_user"
	sessionContextKey contextKey = "auth_session"
)

type AuthMiddleware struct {
	iam *iamclient.Client
}

func NewAuthMiddleware(iam *iamclient.Client) *AuthMiddleware {
	return &AuthMiddleware{iam: iam}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.Header.Get(iamclient.SessionUUIDHeader)
		if sessionUUID == "" {
			writeJSONError(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		user, err := m.iam.Whoami(r.Context(), sessionUUID)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		ctx = context.WithValue(ctx, sessionContextKey, sessionUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) (*iamclient.User, bool) {
	user, ok := ctx.Value(userContextKey).(*iamclient.User)
	return user, ok
}

func SessionFromContext(ctx context.Context) (string, bool) {
	session, ok := ctx.Value(sessionContextKey).(string)
	return session, ok
}

func UserUUIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	user, ok := UserFromContext(ctx)
	if !ok || user == nil {
		return uuid.Nil, false
	}
	return user.UserUUID, true
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    status,
		"message": message,
		"error":   code,
	})
}
