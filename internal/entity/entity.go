package entity

import (
	"errors"
	"time"
)

// User - структура пользователя
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// TeamMember - участник команды
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// Team - команда
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// PullRequest - полная информация о pr
type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`             // OPEN / MERGED
	AssignedReviewers []string   `json:"assigned_reviewers"` // хранится в отдельной таблице pr_reviewers
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

// PullRequestShort - сокращенный pr
type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"` // OPEN / MERGED
}

// ErrorResponse - ответ ошибки
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail - информация об ошибке
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Ошибки
var (
	ErrTeamExists = errors.New("team already exists")
	ErrNotFound   = errors.New("user not found")
	ErrPrExists   = errors.New("PR id already exists")
	ErrPrMerged   = errors.New("cannot reassign on merged PR")
)
