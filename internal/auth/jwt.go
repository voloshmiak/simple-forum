package auth

type jwtAuthenticator struct {
	secret      string
	expiryHours int
}
