package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required, min=7, max=20"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required, min=7, max=20"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GenerateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func RegisterAuthRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	r.POST("/register", Register(pool))
	r.POST("/login", Login(pool))
}

func Register(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignupRequest
		ctx := c.Request.Context()

		// Parse HTTP req, handle errors
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := req.Username
		email := req.Email
		password := req.Password

		hashedPassword, err := HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Push into DB, handling duplicate/invalid errors
		_, err = pool.Exec(ctx,
			"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
			username, email, hashedPassword)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "username is already taken"})
				return
			}
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}

func Login(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse HTTP request
		var req LoginRequest
		ctx := c.Request.Context()

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		username := req.Username
		password := req.Password

		// Extract the salt from the password to compare
		var storedHash string
		var UserID int64
		err := pool.QueryRow(ctx,
			"SELECT id, password_hash FROM users WHERE username=$1",
			username,
		).Scan(&UserID, &storedHash)

		// Handle error s.t. no username match was found
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}

		// Compare the submitted password and the stored hash (password)
		err = CheckPasswordHash(password, storedHash)

		if err != nil {
			// Passwords did not match
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		}

		// Otherwise, generate and send the JWT token
		jwtToken, err := GenerateJWT(UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate a JWT token"})
		}

		c.JSON(http.StatusOK, gin.H{"token": jwtToken})
	}
}
