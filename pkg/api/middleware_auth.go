package api

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/llkhacquan/cisab/pkg/authctx"
	"github.com/llkhacquan/cisab/pkg/models"
	"github.com/llkhacquan/cisab/pkg/repo"
	"github.com/llkhacquan/cisab/pkg/utils/logger"
	"github.com/pkg/errors"
)

// AuthMiddleware validates JWT tokens and sets user information in the request context
func AuthMiddleware(log *logger.Logger, userRepo repo.UserRepo, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for certain public endpoints
			if isPublicEndpoint(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from Authorization header
			tokenString, err := extractTokenFromHeader(r)
			if err != nil {
				log.Info("authorization header invalid", "error", err.Error(), "path", r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			// Parse and validate the token
			token, err := parseAndValidateToken(tokenString, jwtSecret)
			if err != nil {
				log.Info("invalid JWT token", "error", err.Error(), "path", r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// Extract user ID from token
			userID, err := extractUserIDFromToken(token)
			if err != nil {
				log.Info("failed to extract user ID from token", "error", err.Error(), "path", r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			// Fetch the user from the database to get the latest user data
			user, err := userRepo.GetUserByID(r.Context(), userID)
			if err != nil {
				log.Error("failed to fetch user", "error", err.Error(), "user_id", userID, "path", r.URL.Path)
				respondWithError(w, http.StatusInternalServerError, "server error")
				return
			}

			if user == nil {
				log.Info("user not found", "user_id", userID, "path", r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "user not found")
				return
			}

			// Create auth metadata and set it in the context
			auth := authctx.AuthMD{
				User: *user,
			}

			// Set auth context and continue
			ctx := authctx.Set(r.Context(), auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isPublicEndpoint determines if the path is a public endpoint that doesn't require authentication
func isPublicEndpoint(path string) bool {
	// Health check is outside the API router, so we don't need to include it here
	publicPaths := []string{
		"/login", // Login endpoint
		"/users", // Allow creating new users without authentication
	}

	for _, pp := range publicPaths {
		// Match paths that end with the publicPath
		if strings.HasSuffix(path, pp) {
			return true
		}
	}
	return false
}

// respondWithError writes an error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}

// extractUserIDFromToken extracts the user ID from JWT token claims
func extractUserIDFromToken(token *jwt.Token) (models.UserID, error) {
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("failed to extract JWT claims")
	}

	// Extract user ID from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id claim missing or invalid")
	}

	return models.UserID(int(userIDFloat)), nil
}

// parseAndValidateToken parses and validates a JWT token
func parseAndValidateToken(tokenString, jwtSecret string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return token, nil
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(r *http.Request) (string, error) {
	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	// Check if the Authorization header has the correct format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
