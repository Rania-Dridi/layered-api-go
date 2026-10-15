package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rania-Dridi/layered-api-go/internal/handler"
	"github.com/Rania-Dridi/layered-api-go/internal/model"
	"github.com/Rania-Dridi/layered-api-go/internal/repository"
	"github.com/Rania-Dridi/layered-api-go/internal/service"
)

func newServer() http.Handler {
	repo := repository.NewMemory([]model.Course{{ID: 1, Title: "SQL", Capacity: 1}})
	return handler.Router(service.New(repo))
}

func do(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestEnrollThenGrade(t *testing.T) {
	h := newServer()
	if rec := do(h, "POST", "/api/enrollments", `{"courseId":1,"studentId":10}`); rec.Code != http.StatusCreated {
		t.Fatalf("enroll: got %d", rec.Code)
	}
	if rec := do(h, "PUT", "/api/enrollments/1/grade", `{"grade":15.5}`); rec.Code != http.StatusOK {
		t.Fatalf("grade: got %d", rec.Code)
	}
}

func TestBusinessRules(t *testing.T) {
	h := newServer()
	do(h, "POST", "/api/enrollments", `{"courseId":1,"studentId":10}`)
	cases := []struct {
		name, method, path, body string
		want                     int
	}{
		{"course full", "POST", "/api/enrollments", `{"courseId":1,"studentId":11}`, http.StatusConflict},
		{"unknown course", "POST", "/api/enrollments", `{"courseId":99,"studentId":11}`, http.StatusNotFound},
		{"grade out of range", "PUT", "/api/enrollments/1/grade", `{"grade":25}`, http.StatusUnprocessableEntity},
		{"bad json", "POST", "/api/enrollments", `{`, http.StatusBadRequest},
	}
	for _, c := range cases {
		if rec := do(h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s: got %d, want %d", c.name, rec.Code, c.want)
		}
	}
}
