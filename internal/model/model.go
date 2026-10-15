package model

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrCourseFull    = errors.New("course is full")
	ErrAlreadyEnroll = errors.New("student already enrolled")
	ErrInvalidGrade  = errors.New("grade must be between 0 and 20")
)

type Course struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Capacity int    `json:"capacity"`
}

type Enrollment struct {
	ID        int      `json:"id"`
	CourseID  int      `json:"courseId"`
	StudentID int      `json:"studentId"`
	Grade     *float64 `json:"grade,omitempty"`
}
