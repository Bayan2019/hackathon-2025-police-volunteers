package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/repositories/database"
	"github.com/Bayan2019/hackathon-2025-police-volunteers/views"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	// TokenTypeAccess -
	// Set the Issuer to "ozinshe"
	TokenTypeAccess TokenType = "hackathon-2025-access"
)

type AuthHandlers struct {
	DB        *database.Queries
	JwtSecret string
}

func NewAuthHandlers(db *database.Queries, jwtSecret string) *AuthHandlers {
	return &AuthHandlers{
		DB:        db,
		JwtSecret: jwtSecret,
	}
}

type authedHandler func(http.ResponseWriter, *http.Request, views.User)

func (ah *AuthHandlers) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := getBearerToken(r.Header)
		if err != nil {
			views.RespondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
			return
		}

		phone, err := validateJWT(jwtToken, ah.JwtSecret)
		if err != nil {
			views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get email from token", err)
			return
		}

		user, err := ah.DB.GetUserByPhone(r.Context(), phone)
		if err != nil {
			views.RespondWithError(w, http.StatusNotFound, "Couldn't get user", err)
			return
		}

		_, err = ah.DB.GetRefreshTokenOfUser(r.Context(), user.ID)
		if err != nil {
			views.RespondWithError(w, http.StatusInternalServerError, "Couldn't find refresh_token for user", err)
			return
		}

		roles, err := ah.DB.GetRolesOfUser(r.Context(), user.ID)
		if err != nil {
			views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get roles", err)
			return
		}

		user_final := views.User{
			Id:   user.ID,
			Name: user.Name,
			// Phone:       user.Email,
			DateOfBirth: user.DateOfBirth,
			Phone:       user.Phone,
			Roles:       roles,
			Iin:         user.Iin,
		}
		handler(w, r, user_final)
	}
}

// SignIn godoc
// @Tags Auth
// @Summary      Sign In
// @Accept       json
// @Produce      json
// @Param request body views.SignInRequest true "Authentication"
// @Success      200  {object} views.TokensResponse "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid Data"
// @Failure   	 401  {object} views.ErrorResponse "Incorrect email or password"
// @Failure      404  {object} views.ErrorResponse "Email not found"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't create tokens"
// @Router       /v1/auth/sign-in [post]
func (ah *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var signInReq views.SignInRequest
	err := decoder.Decode(&signInReq)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid Data", err)
		return
	}

	user, err := ah.DB.GetUserByPhone(r.Context(), signInReq.Phone)
	if err != nil {
		views.RespondWithError(w, http.StatusNotFound, "Couldn't find the user with such email", err)
		return
	}

	err = checkPasswordHash(signInReq.Password, user.PasswordHash)
	if err != nil {
		views.RespondWithError(w, http.StatusUnauthorized, "Incorrect password", err)
		return
	}

	accessToken, err := makeJWT(
		user.Phone,
		ah.JwtSecret,
		time.Hour*24,
	)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := makeRefreshToken()
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	err = ah.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	})
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token in DataBase", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Refresh godoc
// @Tags Auth
// @Summary      Refresh
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer RefreshToken"
// @Success      200  {object} views.TokensResponse
// @Failure      400  {object} views.ErrorResponse "Couldn't find token"
// @Failure   	 404  {object} views.ErrorResponse "Couldn't find user"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't create tokens"
// @Router       /v1/auth/refresh [post]
func (ah *AuthHandlers) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := getBearerToken(r.Header)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	user, err := ah.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		views.RespondWithError(w, http.StatusNotFound, "Couldn't get user of refresh token", err)
		return
	}

	accessToken, err := makeJWT(
		user.Phone,
		ah.JwtSecret,
		time.Hour,
	)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// SignOut godoc
// @Tags Auth
// @Summary      Sign Out
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      204
// @Failure   	 400  {object} views.ErrorResponse "Couldn't find token"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't revoke session"
// @Router       /v1/auth/sign-out [post]
func (ah *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := getBearerToken(r.Header)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	err = ah.DB.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getBearerToken(headers http.Header) (string, error) {
	// Auth information will come into our server
	// in the Authorization header.
	authHeader := headers.Get("Authorization")
	// fmt.Println(authHeader)
	if authHeader == "" {
		// If the header doesn't exist, return an error.
		return "", errors.New("no auth header included in request")
	}
	// stripping off the Bearer prefix and whitespace
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		// If the header doesn't exist, return an error.
		return "", errors.New("malformed authorization header")
	}
	// return the TOKEN_STRING if it exists
	return splitAuth[1], nil
	// return authHeader, nil
}

func makeJWT(
	email string,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	// Use jwt.NewWithClaims to create a new token
	token := jwt.NewWithClaims(
		// Use jwt.SigningMethodHS256 as the signing method.
		jwt.SigningMethodHS256,
		// Use jwt.RegisteredClaims as the claims
		jwt.RegisteredClaims{
			Issuer: string(TokenTypeAccess),
			// Set IssuedAt to the current time in UTC
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			// Set ExpiresAt to the current time plus the expiration time (expiresIn)
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			// Set the Subject to a stringified version of the user's email
			Subject: email,
		})
	// Use token.SignedString to sign the token with the secret key.
	return token.SignedString(signingKey)
}

func validateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	// Use the jwt.ParseWithClaims function
	// to validate the signature of the JWT
	// and extract the claims into a *jwt.Token struct.
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return "", err
	}

	// If all is well with the token,
	// use the token.Claims interface
	// to get access to the user's email from the claims
	// (which should be stored in the Subject field).
	email, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeAccess) {
		return "", errors.New("invalid issuer")
	}

	return email, nil
}

// Hash the password using the bcrypt.GenerateFromPassword function
// HashPassword -
func hashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func makeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
