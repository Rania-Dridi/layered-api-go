// Package service holds the business rules (capacity, duplicates, grading).
package service

import (
	"github.com/Rania-Dridi/layered-api-go/internal/model"
	"github.com/Rania-Dridi/layered-api-go/internal/repository"
)

type Service struct{ repo repository.Repository }

func New(r repository.Repository) *Service { return &Service{repo: r} }

func (s *Service) Courses() []model.Course { return s.repo.ListCourses() }

func (s *Service) Enroll(courseID, studentID int) (model.Enrollment, error) {
	course, err := s.repo.GetCourse(courseID)
	if err != nil {
		return model.Enrollment{}, err
	}
	existing := s.repo.ListEnrollments(courseID)
	if len(existing) >= course.Capacity {
		return model.Enrollment{}, model.ErrCourseFull
	}
	for _, e := range existing {
		if e.StudentID == studentID {
			return model.Enrollment{}, model.ErrAlreadyEnroll
		}
	}
	return s.repo.AddEnrollment(model.Enrollment{CourseID: courseID, StudentID: studentID}), nil
}

func (s *Service) Grade(enrollmentID int, grade float64) (model.Enrollment, error) {
	if grade < 0 || grade > 20 {
		return model.Enrollment{}, model.ErrInvalidGrade
	}
	e, err := s.repo.GetEnrollment(enrollmentID)
	if err != nil {
		return model.Enrollment{}, err
	}
	e.Grade = &grade
	if err := s.repo.UpdateEnrollment(e); err != nil {
		return model.Enrollment{}, err
	}
	return e, nil
}
