package dexp

import "github.com/dgrijalva/jwt-go"

var mySigningKey = []byte("dexp.io@2019")

func JwtDecode(token string) (*jwt.Token, error) {

	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

}

func JwtCreate(userID int64, expiredAt int64) string {

	claims := UserClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    "dexp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)

	return ss
}
