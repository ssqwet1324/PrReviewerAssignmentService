package repository

import (
	"context"
	"errors"
	"fmt"
	"pr_reviewer_service/internal/config"
	"pr_reviewer_service/internal/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Repository - бд
type Repository struct {
	DB     *pgxpool.Pool
	Logger *zap.Logger
}

// New - конструктор для repo
func New(cfg *config.Config, logger *zap.Logger) (*Repository, error) {
	dsn := cfg.GetDSN()

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("repository: Error connecrtion from pgxpool: %v", err)
	}

	return &Repository{
		DB:     dbPool,
		Logger: logger,
	}, nil
}

// CreateTeam - создать команду
func (repo *Repository) CreateTeam(ctx context.Context, team entity.Team) error {
	tx, err := repo.DB.Begin(ctx)
	if err != nil {
		repo.Logger.Error("Error begin transaction", zap.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				repo.Logger.Error("Error rollback", zap.Error(rbErr))
			}
			return
		}

		if cmErr := tx.Commit(ctx); cmErr != nil {
			repo.Logger.Error("Error commit", zap.Error(cmErr))
			err = cmErr
		}
	}()

	_, err = tx.Exec(ctx, `INSERT INTO team (team_name) VALUES ($1)`, team.TeamName)
	if err != nil {
		repo.Logger.Error("Error insert into team", zap.Error(err))
		return err
	}

	for _, member := range team.Members {
		_, err := tx.Exec(ctx, `INSERT INTO users (team_name, user_id, username, is_active) VALUES ($1, $2, $3, $4)`,
			team.TeamName, member.UserID, member.Username, member.IsActive)
		if err != nil {
			repo.Logger.Error("Error insert into team_member", zap.Error(err))
			return err
		}
	}

	return nil
}

// GetTeam - получить команду и ее участников
func (repo *Repository) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	var team entity.Team
	var members []entity.TeamMember

	team.TeamName = teamName

	rows, err := repo.DB.Query(ctx, `SELECT user_id, username, is_active FROM users
    	WHERE team_name = $1`, teamName)
	if err != nil {
		repo.Logger.Error("Error select from team", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m entity.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			repo.Logger.Error("Error scan from team", zap.Error(err))
			return nil, err
		}
		members = append(members, m)
	}

	team.Members = members

	return &team, nil
}

// ChangeActivityUser - изменить активность пользователя
func (repo *Repository) ChangeActivityUser(ctx context.Context, isActive bool, userID string) error {
	_, err := repo.DB.Exec(ctx, `UPDATE users SET is_active = $1 WHERE user_id = $2`,
		isActive, userID)
	if err != nil {
		repo.Logger.Error("Error update user", zap.Error(err))
		return err
	}

	return nil
}

// GetReviewFromUser - получить pr где пользователь ревьювер
func (repo *Repository) GetReviewFromUser(ctx context.Context, userID string) ([]entity.PullRequestShort, error) {
	var prs []entity.PullRequestShort

	rows, err := repo.DB.Query(ctx, `SELECT p.pull_request_id, p.pull_request_name, p.author_id, p.status
		FROM pr_reviewers r
		JOIN pr p ON r.pull_request_id = p.pull_request_id
		WHERE r.user_id = $1`, userID)
	if err != nil {
		repo.Logger.Error("Error selecting PRs for user", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pr entity.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			repo.Logger.Error("Error scanning PR", zap.Error(err))
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

// CreatePullRequest - создать новый pr
func (repo *Repository) CreatePullRequest(ctx context.Context, pr entity.PullRequest) error {
	tx, err := repo.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("CreatePullRequest: begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				repo.Logger.Error("Error rollback", zap.Error(rbErr))
			}
			return
		}

		if cmErr := tx.Commit(ctx); cmErr != nil {
			repo.Logger.Error("Error commit", zap.Error(cmErr))
			err = cmErr
		}
	}()

	_, err = tx.Exec(ctx, `
        INSERT INTO pr (pull_request_id, pull_request_name, author_id, status, created_at)
        VALUES ($1, $2, $3, $4, $5)`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt)
	if err != nil {
		repo.Logger.Error("CreatePullRequest: Failed to insert PR", zap.Error(err))
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err = tx.Exec(ctx, `
            INSERT INTO pr_reviewers (pull_request_id, user_id)
            VALUES ($1, $2)`,
			pr.PullRequestID, reviewerID)
		if err != nil {
			repo.Logger.Error("CreatePullRequest: Failed to insert reviewer", zap.Error(err), zap.String("reviewer_id", reviewerID))
			return err
		}
	}

	repo.Logger.Info("Pull request created", zap.String("pr_id", pr.PullRequestID))

	return nil
}

// GetPR - получить pr
func (repo *Repository) GetPR(ctx context.Context, pullRequestID string) (entity.PullRequest, error) {
	var pr entity.PullRequest

	row := repo.DB.QueryRow(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pr WHERE pull_request_id = $1`, pullRequestID)

	var mergedAt *time.Time
	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	if err != nil {
		repo.Logger.Error("Error selecting PR", zap.Error(err))
		return pr, entity.ErrNotFound
	}
	pr.MergedAt = mergedAt

	rows, err := repo.DB.Query(ctx, `SELECT user_id FROM pr_reviewers WHERE pull_request_id = $1`,
		pullRequestID)
	if err != nil {
		repo.Logger.Error("Error selecting PR reviewers", zap.Error(err))
		return pr, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			repo.Logger.Error("Error scanning reviewer", zap.Error(err))
			return pr, err
		}
		reviewers = append(reviewers, userID)
	}

	if rows.Err() != nil {
		repo.Logger.Error("Error selecting PR reviewers", zap.Error(rows.Err()))
		return pr, rows.Err()
	}

	pr.AssignedReviewers = reviewers

	return pr, nil
}

// UpdatePRStatus - обновить статус pr
func (repo *Repository) UpdatePRStatus(ctx context.Context, prID, newPrStatus string) error {
	_, err := repo.DB.Exec(ctx, `UPDATE pr SET status = $1 WHERE pull_request_id = $2`, newPrStatus, prID)
	if err != nil {
		repo.Logger.Error("Error update PR status", zap.Error(err))
		return err
	}

	return nil
}

// MergePr - меняем статус pr
func (repo *Repository) MergePr(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr, err := repo.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == "MERGED" {
		return &pr, nil
	}

	now := time.Now()
	pr.Status = "MERGED"
	pr.MergedAt = &now

	// Обновляем и статус, и merged_at в БД
	_, err = repo.DB.Exec(ctx, `UPDATE pr SET status = $1, merged_at = $2 WHERE pull_request_id = $3`,
		pr.Status, now, prID)
	if err != nil {
		repo.Logger.Error("Error update PR status and merged_at", zap.Error(err))
		return nil, err
	}

	return &pr, nil
}

// ReassignPrReviewer - переназначить ревьюера
func (repo *Repository) ReassignPrReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (entity.PullRequest, error) {
	pr, err := repo.GetPR(ctx, prID)
	if err != nil {
		return pr, err
	}

	// Проверяем, что PR открыт
	if pr.Status == "MERGED" {
		return pr, fmt.Errorf("cannot reassign reviewer on merged PR")
	}

	// Проверяем, что oldReviewer действительно назначен
	found := false
	for _, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return pr, fmt.Errorf("reviewer %s is not assigned to this PR", oldReviewerID)
	}

	// Переназначаем в таблице pr_reviewers
	cmdTag, err := repo.DB.Exec(ctx, `UPDATE pr_reviewers SET user_id = $1 WHERE pull_request_id = $2
        AND user_id = $3`, newReviewerID, prID, oldReviewerID)
	if err != nil {
		repo.Logger.Error("Error reassign PR reviewer", zap.Error(err))
		return pr, err
	}

	if cmdTag.RowsAffected() == 0 {
		return pr, fmt.Errorf("no reviewer updated")
	}

	// Обновляем структуру PR в памяти
	for i, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewerID
			break
		}
	}

	repo.Logger.Info("PR reviewer reassigned",
		zap.String("pr_id", prID),
		zap.String("old_reviewer", oldReviewerID),
		zap.String("new_reviewer", newReviewerID),
	)

	return pr, nil
}

// GetTeamByUserID - получить имя команды по id пользователя
func (repo *Repository) GetTeamByUserID(ctx context.Context, userID string) (string, error) {
	var teamName string

	err := repo.DB.QueryRow(ctx, `SELECT team_name FROM users WHERE user_id=$1`, userID).Scan(&teamName)
	if err != nil {
		repo.Logger.Error("Error getting team by user ID", zap.Error(err))
		return "", entity.ErrNotFound
	}

	return teamName, nil
}

// CheckTeam - проверка на существование команды
func (repo *Repository) CheckTeam(ctx context.Context, teamName string) (bool, error) {
	var exists int

	//err := repo.DB.QueryRow(ctx, `SELECT 1 FROM users WHERE team_name=$1 LIMIT 1`, teamName).Scan(&exists)

	err := repo.DB.QueryRow(ctx, `SELECT 1 FROM team WHERE team_name=$1 LIMIT 1`, teamName).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		repo.Logger.Error("Error checking team", zap.Error(err))
		return false, err
	}

	return true, nil
}

// CheckUser - проверяем существует ли пользователь
func (repo *Repository) CheckUser(ctx context.Context, userID string) (bool, error) {
	var exists int

	err := repo.DB.QueryRow(ctx, `SELECT 1 FROM users WHERE user_id=$1`, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		repo.Logger.Error("Error checking user", zap.Error(err))
		return false, err
	}

	return true, nil
}

// CheckPR - проверить существует ли pr
func (repo *Repository) CheckPR(ctx context.Context, prID string) (bool, error) {
	var exists int

	err := repo.DB.QueryRow(ctx, `SELECT 1 FROM pr WHERE pull_request_id=$1 LIMIT 1`, prID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		repo.Logger.Error("Error checking PR", zap.Error(err))
		return false, err
	}

	return true, nil
}
