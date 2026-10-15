// Package repository isolates data access. The in-memory implementation can be
// swapped for a PostgreSQL one without touching the service layer.
package repository

import (
	"sync"

	"github.com/Rania-Dridi/layered-api-go/internal/model"
)

type Repository interface {
	ListCourses() []model.Course
	GetCourse(id int) (model.Course, error)
	ListEnrollments(courseID int) []model.Enrollment
	AddEnrollment(e model.Enrollment) model.Enrollment
	GetEnrollment(id int) (model.Enrollment, error)
	UpdateEnrollment(e model.Enrollment) error
}

type Memory struct {
	mu          sync.RWMutex
	courses     []model.Course
	enrollments []model.Enrollment
	nextEnroll  int
}

func NewMemory(courses []model.Course) *Memory {
	return &Memory{courses: courses, nextEnroll: 1}
}

func (m *Memory) ListCourses() []model.Course {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]model.Course(nil), m.courses...)
}

func (m *Memory) GetCourse(id int) (model.Course, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.courses {
		if c.ID == id {
			return c, nil
		}
	}
	return model.Course{}, model.ErrNotFound
}

func (m *Memory) ListEnrollments(courseID int) []model.Enrollment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []model.Enrollment
	for _, e := range m.enrollments {
		if e.CourseID == courseID {
			out = append(out, e)
		}
	}
	return out
}

func (m *Memory) AddEnrollment(e model.Enrollment) model.Enrollment {
	m.mu.Lock()
	defer m.mu.Unlock()
	e.ID = m.nextEnroll
	m.nextEnroll++
	m.enrollments = append(m.enrollments, e)
	return e
}

func (m *Memory) GetEnrollment(id int) (model.Enrollment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, e := range m.enrollments {
		if e.ID == id {
			return e, nil
		}
	}
	return model.Enrollment{}, model.ErrNotFound
}

func (m *Memory) UpdateEnrollment(e model.Enrollment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.enrollments {
		if m.enrollments[i].ID == e.ID {
			m.enrollments[i] = e
			return nil
		}
	}
	return model.ErrNotFound
}
