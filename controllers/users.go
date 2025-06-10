package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/repositories"
	"github.com/Bayan2019/hackathon-2025-police-volunteers/repositories/database"
	"github.com/Bayan2019/hackathon-2025-police-volunteers/views"
	"github.com/go-chi/chi"
)

type UsersHandlers struct {
	userRepo *repositories.UsersRepository
}

func NewUsersHandlers(repo *repositories.UsersRepository) *UsersHandlers {
	return &UsersHandlers{
		userRepo: repo,
	}
}

// Register godoc
// @Tags Users
// @Summary      Create user (Register)
// @Accept       json
// @Produce      json
// @Param request body views.CreateUserRequest true "User data"
// @Success      201  {object} views.ResponseId "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid data"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't hash password"
// @Router       /v1/users [post]
func (uh *UsersHandlers) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	cur := views.CreateUserRequest{}

	err := decoder.Decode(&cur)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Error parsing JSON of CreateUserRequest", err)
		return
	}

	hashedPassword, err := hashPassword(cur.Password)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	id, err := uh.userRepo.Create(r.Context(), database.CreateUserParams{
		Name:            cur.Name,
		Phone:           cur.Phone,
		Iin:             cur.Iin,
		DateOfBirth:     cur.DateOfBirth,
		CurrentLocation: cur.CurrentLocation,
		PasswordHash:    hashedPassword,
	})
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	views.RespondWithJSON(w, http.StatusCreated, views.NewResponseId(int(id)))
}

// UpdateProfile godoc
// @Tags Users
// @Summary      Update user profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Param request body views.UpdateProfileRequest true "User data"
// @Success      200  "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid data"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't update user data"
// @Router       /v1/users/profile [put]
// @Security Bearer
func (uh *UsersHandlers) UpdateProfile(w http.ResponseWriter, r *http.Request, user views.User) {
	decoder := json.NewDecoder(r.Body)
	upr := views.UpdateProfileRequest{}

	err := decoder.Decode(&upr)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Error parsing JSON of UpdateProfileRequest", err)
		return
	}

	err = uh.userRepo.UpdateProfile(r.Context(), user.Id, upr)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user data", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetProfile godoc
// @Tags Users
// @Summary      Get User profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {object} views.User "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/users/profile [get]
// @Security Bearer
func (uh *UsersHandlers) GetProfile(w http.ResponseWriter, r *http.Request, user views.User) {
	views.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteProfile godoc
// @Tags Users
// @Summary      Delete user profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {object} views.ResponseId "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't delete user"
// @Router       /v1/users/profile [delete]
// @Security Bearer
func (uh *UsersHandlers) DeleteProfile(w http.ResponseWriter, r *http.Request, user views.User) {
	err := uh.userRepo.DB.DeleteUser(r.Context(), user.Id)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete user", err)
		return
	}
	views.RespondWithJSON(w, http.StatusOK, views.NewResponseId(int(user.Id)))
}

// Update godoc
// @Tags Users
// @Summary      Update user
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Param id path int true "id"
// @Param request body views.UpdateUserRequest true "User data"
// @Success      200  "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid data"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 403  {object} views.ErrorResponse "No Permission"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't update user data"
// @Router       /v1/users/{id} [put]
// @Security Bearer
func (uh *UsersHandlers) Update(w http.ResponseWriter, r *http.Request, user views.User) {

	can_do := false
	for _, role := range user.Roles {
		if role.Title == "admin" {
			can_do = true
			break
		}
	}

	if !can_do {
		views.RespondWithError(w, http.StatusForbidden, "don't have permission", errors.New("no Permission"))
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid id", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	uur := views.UpdateUserRequest{}

	err = decoder.Decode(&uur)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Error parsing JSON of UpdateUserRequest", err)
		return
	}

	err = uh.userRepo.Update(r.Context(), int64(id), uur)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user data", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete godoc
// @Tags Users
// @Summary      Get User
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Param id path int true "id"
// @Success      200  {object} views.User "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid data"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 403  {object} views.ErrorResponse "No Permission"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't get user"
// @Router       /v1/users/{id} [get]
// @Security Bearer
func (uh *UsersHandlers) GetUser(w http.ResponseWriter, r *http.Request, user views.User) {
	can_do := false
	for _, role := range user.Roles {
		if role.Title == "admin" {
			can_do = true
			break
		}
	}
	if !can_do {
		views.RespondWithError(w, http.StatusForbidden, "don't have permission", errors.New("no Permission"))
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid id", err)
		return
	}

	user1, err := uh.userRepo.DB.GetUserById(r.Context(), int64(id))
	if err != nil {
		views.RespondWithError(w, http.StatusNotFound, "Couldn't get user", err)
		return
	}

	roles, err := uh.userRepo.DB.GetRolesOfUser(r.Context(), user1.ID)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get roles", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.User{
		Id:   user1.ID,
		Name: user1.Name,
		// Email:       user1.Email,
		DateOfBirth: user1.DateOfBirth,
		Phone:       user1.Phone,
		Iin:         user1.Iin,
		Roles:       roles,
	})
}

// GetUsers godoc
// @Tags Users
// @Summary      Get Users List
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {array} views.User "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 403  {object} views.ErrorResponse "No Permission"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get users"
// @Router       /v1/users [get]
// @Security Bearer
func (uh *UsersHandlers) GetUsers(w http.ResponseWriter, r *http.Request, user views.User) {
	can_do := false
	for _, role := range user.Roles {
		if role.Title == "admin" {
			can_do = true
			break
		}
	}
	if !can_do {
		views.RespondWithError(w, http.StatusForbidden, "don't have permission", errors.New("no Permission"))
		return
	}

	users, err := uh.userRepo.DB.GetUsers(r.Context())
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get users", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, users)
}

// Delete godoc
// @Tags Users
// @Summary      Delete user profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Param id path int true "id"
// @Success      200  {object} views.ResponseId "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid data"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 403  {object} views.ErrorResponse "No Permission"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't delete user"
// @Router       /v1/users/{id} [delete]
// @Security Bearer
func (uh *UsersHandlers) Delete(w http.ResponseWriter, r *http.Request, user views.User) {
	can_do := false
	for _, role := range user.Roles {
		if role.Title == "admin" {
			can_do = true
			break
		}
	}

	if !can_do {
		views.RespondWithError(w, http.StatusForbidden, "don't have permission", errors.New("no Permission"))
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid id", err)
		return
	}

	err = uh.userRepo.DB.DeleteUser(r.Context(), int64(id))
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete user", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.NewResponseId(id))
}
