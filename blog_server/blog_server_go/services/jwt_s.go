package services

/*
厦门大学计算机专业 | 前华为工程师
专注《零基础学编程系列》  http://lblbc.cn/blog
包含：Java | 安卓 | 前端 | Flutter | iOS | 小程序 | 鸿蒙
公众号：蓝不蓝编程
*/
import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(userID string) string             // Generate a new token
	ValidateToken(token string) (*jwt.Token, error) // Validate the token
}

// jwtCustomClaims is a struct that contains the custom claims for the JWT
type jwtCustomClaim struct {
	UserID             string `json:"user_id"` // The userId is the only required field
	jwt.StandardClaims        // This is a standard JWT claim
}

// jwtService is a struct that implements the JWTService interface
type jwtService struct {
	secretKey string // Secret key used to sign the token
	issuer    string // Who creates the token
}

//NewJWTService method is creates a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "gojwt",        // who creates the token
		secretKey: getSecretKey(), // Call the getSecretKey function to get the secret key
	}
}

// Create get the secret key from the environment variable
func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET_KEY") // Get the secret key from the environment variable
	if secretKey == "" {
		secretKey = "secret" // If the environment variable is empty, use a default value
	}
	return secretKey
}

// Create a new token object, specifying signing method and the claims
func (s *jwtService) GenerateToken(userID string) string {

	// Create the Claims struct with the required claims for the JWT
	claims := &jwtCustomClaim{
		userID, // userId is the only required field
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(), // 1 year expiration
			Issuer:    s.issuer,                           // Who creates the token
			IssuedAt:  time.Now().Unix(),                  // when the token was issued/created (now)
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Sign the token with our secret
	t, err := token.SignedString([]byte(s.secretKey))          // Sign the token with an expiration time
	if err != nil {
		panic(err) // If there is an error, panic
	}
	return t // Return the token to the user, along with an expiration time
}

// ValidateToken validates the token and returns the claims
func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	// Parse the token
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok { // Check the signing method
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"]) // Return an error if the signing method isn't HMAC
		}
		return []byte(s.secretKey), nil // Return the key
	})
}
