package easyAuth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/khengsaurus/easy-auth-go/consts"
	"github.com/khengsaurus/easy-auth-go/types"
)

type validateRes struct {
	Data types.User
}

var (
	client = http.Client{Timeout: 5 * time.Second}
)

func validateToken(token string, version int, apiKey string) (*types.User, error) {
	url := fmt.Sprintf("%s/v%d/ea_validate/%s", consts.UrlBase, version, token)
	ssoReq, err := http.NewRequest(http.MethodGet, url, nil)
	ssoReq.Header.Set("Content-Type", "application/json")
	ssoReq.Header.Set(consts.HeaderApiKey, apiKey)
	if err != nil {
		return nil, err
	}

	ssoRes, err := client.Do(ssoReq)
	if err != nil {
		return nil, err
	}
	defer ssoRes.Body.Close()

	if ssoRes.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf(consts.ErrorInvalidToken)
	}

	body, err := io.ReadAll(ssoRes.Body)
	if err != nil {
		return nil, err
	}

	var p validateRes
	err = json.Unmarshal(body, &p)
	return &p.Data, err
}

// Retrieves user info from request header "Ea-User-Token" and stores it in request Context.
//
// Param silent=true: if header value is missing, invalid, or fails to retrieve user info, will not raise an error.
//
// Param silent=false: if header value is missing or invalid, will write a 401 response. If user info cannot be retrieved for other reasons, will write a 500 response.
func WithEaUser(apiKey string, version int, silent bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(consts.HeaderUserToken)
			if token == "" {
				if silent {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
				}
				return
			}

			sessionUser, err := validateToken(token, version, apiKey)
			if err != nil || sessionUser == nil || sessionUser.Id == "" || sessionUser.Username == "" {
				if silent {
					next.ServeHTTP(w, r)
				} else {
					if err != nil && err.Error() == consts.ErrorInvalidToken {
						w.WriteHeader(http.StatusUnauthorized)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
				return
			} else {
				ctx := context.WithValue(r.Context(), consts.ContextKeyUser, sessionUser)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func GetEaUser(ctx context.Context) (*types.User, error) {
	user, ok := ctx.Value(consts.ContextKeyUser).(*types.User)
	if !ok || user.Id == "" || user.Username == "" {
		return nil, fmt.Errorf(consts.ErrorUserInfoMissing)
	}
	return user, nil
}
