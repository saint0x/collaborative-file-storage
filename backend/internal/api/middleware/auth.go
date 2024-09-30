package apimiddleware

import (
	"net/http"

	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/pkg/errors"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func Auth(authService *auth.ClerkService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.ExtractBearerToken(r)
			if err != nil {
				utils.RespondError(w, errors.Unauthorized("Invalid or missing token"))
				return
			}

			userID, err := authService.ValidateAndExtractUserID(r.Context(), token)
			if err != nil {
				utils.RespondError(w, errors.Unauthorized("Invalid token"))
				return
			}

			ctx := auth.SetUserIDInContext(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
