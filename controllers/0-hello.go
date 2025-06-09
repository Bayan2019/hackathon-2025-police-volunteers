package controllers

import (
	"net/http"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/views"
)

// Static godoc
// @Summary      Saying hello
// @Tags Tests
// @Produce      json
// @Success      200  {object} views.ResponseMessage "OK"
// @Failure   	 500  {object} views.ErrorResponse "Invalid file"
// @Router       / [get]
func HelloHandler(w http.ResponseWriter, r *http.Request) {

	views.RespondWithJSON(w, http.StatusOK, views.ResponseMessage{
		Message: "Welcom from Police Volunteers api",
	})
}
