// Package handler translates HTTP requests into service calls and back.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Rania-Dridi/layered-api-go/internal/model"
	"github.com/Rania-Dridi/layered-api-go/internal/service"
)

type Handler struct{ svc *service.Service }

// Router wires the routes to the handlers.
func Router(svc *service.Service) http.Handler {
	h := &Handler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/courses", h.listCourses)
	mux.HandleFunc("POST /api/enrollments", h.enroll)
	mux.HandleFunc("PUT /api/enrollments/{id}/grade", h.grade)
	return mux
}

func (h *Handler) listCourses(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Courses())
}

func (h *Handler) enroll(w http.ResponseWriter, r *http.Request) {
	var in struct {
		CourseID  int `json:"courseId"`
		StudentID int `json:"studentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	e, err := h.svc.Enroll(in.CourseID, in.StudentID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *Handler) grade(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var in struct {
		Grade float64 `json:"grade"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	e, err := h.svc.Grade(id, in.Grade)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, model.ErrCourseFull), errors.Is(err, model.ErrAlreadyEnroll):
		status = http.StatusConflict
	case errors.Is(err, model.ErrInvalidGrade):
		status = http.StatusUnprocessableEntity
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
