claims := &jwt.RegisteredClaims{
	ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	Subject:   fmt.Sprintf("%d", id.Id().ToInt()),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString(a.secretKey)
