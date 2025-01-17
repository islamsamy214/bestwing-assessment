package services

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"
	"web-app/app/models/user"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
)

var (
	// Define a secret key for signing tokens. This should be securely stored in a real-world application.
	secretKey = []byte(os.Getenv("JWT_SECRET"))
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token with the provided user ID
func GenerateToken(userID int64) (string, error) {
	// Create a new set of claims
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
			Issuer:    "github@islamsamy214",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken parses the provided JWT token string and returns the claims if the token is valid
func ParseToken(tokenStr string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrSignatureInvalid
	}
}

// ValidateToken validates the provided JWT token string
func ValidateToken(tokenStr string) error {
	// Parse the token
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	// Generate a salt with a length of 16 bytes
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password using the Argon2id key derivation function
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Encode the salt and hashed password to a base64 string
	hashedPassword := base64.StdEncoding.EncodeToString(append(salt, hash...))

	return hashedPassword, nil
}

func DecodeHashedPassword(hashedPassword string) ([]byte, []byte, error) {
	// Decode the base64 string to the salt and hashed password
	data, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return nil, nil, err
	}

	salt := data[:16]
	hash := data[16:]

	return salt, hash, nil
}

func AttemptLogin(username, password string) (bool, error) {
	// Get the user from the database
	user, err := GetUserByUsername(user.NewUserModel())
	if err != nil {
		return false, err
	}

	// Decode the hashed password
	salt, hashedPassword, err := DecodeHashedPassword(user.Password)
	if err != nil {
		return false, err
	}

	// Hash the provided password with the salt
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Compare the hashed password with the provided password
	if string(hashedPassword) == string(hash) {
		return true, nil
	}

	return false, nil
}

func GetUserByUsername(u *user.User) (*user.User, error) {
	// Find the user by username
	err := u.FindByUsername()
	if err != nil {
		return nil, err
	}

	return u, nil
}
