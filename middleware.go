package dexp

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"net"
	"net/http"
	"strings"
)

type UserAuth struct {
	UserID    int64
	Roles     []string
	IPAddress string
	Token     string
}

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type UserClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := TokenFromHttpRequest(r)
			userId := UserIDFromToken(token)
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			var roles [] string
			userAuth := UserAuth{
				UserID:    userId,
				Roles:     roles,
				IPAddress: ip,
				Token:     token,
			}
			// get the user from the database
			if userId > 0 {
				userAuth.Roles = GetUserRoles(userId)
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &userAuth)
			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserRoles(userID int64) [] string {

	var roles []string
	q, err := DB.Select("user_has_role", "ur").Join("role", "r", "r.id = ur.role_id").Fields("r", []string{"name"}).Condition("ur.user_id", userID, "=").FetchAll()

	if err != nil {
		return roles
	}

	defer q.Close()

	for q.Next() {
		var role string

		if q.Scan(&role) == nil {
			roles = append(roles, role)
		}
	}

	return roles

}

func TokenFromHttpRequest(r *http.Request) string {

	reqToken := r.Header.Get("Authorization")
	var tokenString string

	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}

	return tokenString

}

func UserIDFromToken(tokenString string) int64 {

	if TokenIsBlackList(tokenString) {
		return 0
	}

	token, err := JwtDecode(tokenString)

	if err != nil {
		return 0
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {

		if claims == nil {
			return 0
		}
		return claims.UserId
	} else {
		return 0
	}
}

func ForContext(ctx context.Context) *UserAuth {
	raw := ctx.Value(userCtxKey)

	if raw == nil {
		return nil
	}

	return raw.(*UserAuth)
}

func GetAuthFromContext(ctx context.Context) *UserAuth {
	return ForContext(ctx)
}
