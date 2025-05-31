package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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
	server := echo.New()

	server.POST("/register", func(c echo.Context) error {
		args := new(Login_Request)
		if err := c.Bind(args); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		}

		hashedPassword, err := hashPassword(args.Password)
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Failed to hash password"})
		}

		users[args.Username] = Login{
			HashedPassword: hashedPassword,
		}

		return c.JSON(201, map[string]string{"message": "User registered successfully"})
	})

	server.POST("/login", func(c echo.Context) error {
		args := new(Login_Request)
		if err := c.Bind(args); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
		}

		user, exists := users[args.Username]
		if !exists || !checkPasswordHash(args.Password, user.HashedPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		}

		sessionToken := generateToken(32)
		csrfToken := generateToken(32)

		http.SetCookie(c.Response(), &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(14 * 24 * time.Hour),
			HttpOnly: true,
		})

		http.SetCookie(c.Response(), &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Expires:  time.Now().Add(14 * 24 * time.Hour),
			HttpOnly: false, // CSRF token should be accessible via JavaScript
		})

		user.SessionToken = sessionToken
		user.CSRFToken = csrfToken
		users[args.Username] = user

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Login successful",
		})
	})

	server.GET("/protected", func(c echo.Context) error {
		if err := authorize(c); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Access granted to protected resource"})
	})

	server.POST("/logout", func(c echo.Context) error {
		if err := authorize(c); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		http.SetCookie(c.Response(), &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire the session token
			HttpOnly: true,
		})

		http.SetCookie(c.Response(), &http.Cookie{
			Name:     "csrf_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire the session token
			HttpOnly: false,
		})

		// todo: sollte viellecht eher user id sein
		username := c.Request().Header.Get("X-Username")

		user, exists := users[username]
		if !exists {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		user.SessionToken = ""
		user.CSRFToken = ""
		users[username] = user
		return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
	})

	server.Logger.Fatal(server.Start(":5080"))
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

func authorize(c echo.Context) error {
	// todo: sollte viellecht eher user id sein
	username := c.Request().Header.Get("X-Username")

	user, exists := users[username]
	if !exists {
		return ErrAuth
	}

	sessionToken, err := c.Cookie("session_token")
	if err != nil || sessionToken == nil || sessionToken.Value == "" || sessionToken.Value != user.SessionToken {
		return ErrAuth
	}

	csrfToken := c.Request().Header.Get("X-CSRF-Token")
	if csrfToken == "" || csrfToken != user.CSRFToken {
		return ErrAuth
	}

	return nil
}
