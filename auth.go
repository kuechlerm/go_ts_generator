package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Login_Request struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// todo: lucia nutzt noch eine Session Tabelle; nötig?
type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

// todo: temp DB
var users = map[string]Login{}

func Init_Auth() {
	server := gin.Default()

	server.POST("/register", func(c *gin.Context) {
		args := new(Login_Request)
		if err := c.Bind(args); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		hashedPassword, err := hashPassword(args.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		users[args.Username] = Login{
			HashedPassword: hashedPassword,
		}

		c.JSON(201, gin.H{"message": "User registered successfully"})
	})

	server.POST("/login", func(c *gin.Context) {
		args := new(Login_Request)
		if err := c.Bind(args); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		user, exists := users[args.Username]
		if !exists || !checkPasswordHash(args.Password, user.HashedPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		sessionToken := generateToken(32)
		csrfToken := generateToken(32)

		c.SetCookie(
			"session_token",
			sessionToken,
			// todo: expiration lieber in DB (Login-struct) abbilden
			60*60*24*14, // 14 days in seconds
			"/",
			"localhost", // todo: anpassen für Prod
			false,       // todo: so?
			true,
		)

		c.SetCookie(
			"csrf_token",
			csrfToken,
			// todo: expiration lieber in DB (Login-struct) abbilden
			60*60*24*14, // 14 days in seconds
			"/",
			"localhost", // todo: anpassen für Prod
			false,       // todo: so?
			false,
		)

		user.SessionToken = sessionToken
		user.CSRFToken = csrfToken
		users[args.Username] = user

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
		})
	})

	server.GET("/protected", func(c *gin.Context) {
		if err := authorize(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Access granted to protected resource"})
	})

	server.POST("/logout", func(c *gin.Context) {
		if err := authorize(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.SetCookie(
			"session_token",
			"",
			-1, // delete cookie
			"/",
			"localhost", // todo: anpassen für Prod
			false,       // todo: so?
			true,
		)

		c.SetCookie(
			"csrf_token",
			"",
			// todo: expiration lieber in DB (Login-struct) abbilden
			-1, // 14 days in seconds
			"/",
			"localhost", // todo: anpassen für Prod
			false,       // todo: so?
			false,
		)

		// todo: sollte viellecht eher user id sein
		username := c.GetHeader("X-Username")

		user, exists := users[username]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		user.SessionToken = ""
		user.CSRFToken = ""
		users[username] = user
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	})

	// server.Logger.Fatal(server.Start(":5080"))
	server.Run(":5080")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

var ErrAuth = errors.New("authentication error")

func authorize(c *gin.Context) error {
	// todo: sollte viellecht eher user id sein
	username := c.GetHeader("X-Username")

	user, exists := users[username]
	if !exists {
		return ErrAuth
	}

	sessionToken, err := c.Cookie("session_token")
	if err != nil || sessionToken == "" || sessionToken != user.SessionToken {
		return ErrAuth
	}

	csrfToken := c.GetHeader("X-CSRF-Token")
	if csrfToken == "" || csrfToken != user.CSRFToken {
		return ErrAuth
	}

	return nil
}
