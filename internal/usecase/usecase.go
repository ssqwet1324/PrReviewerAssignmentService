package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"pr_reviewer_service/internal/entity"
	"time"
)

// RepositoryProvider - поведение функций repository
type RepositoryProvider interface {
	CreateTeam(ctx context.Context, team entity.Team) error
	GetTeam(ctx context.Context, teamName string) (*entity.Team, error)
	ChangeActivityUser(ctx context.Context, isActive bool, userID string) error
	GetReviewFromUser(ctx context.Context, userID string) ([]entity.PullRequestShort, error)
	CreatePullRequest(ctx context.Context, pr entity.PullRequest) error
	GetPR(ctx context.Context, pullRequestID string) (entity.PullRequest, error)
	UpdatePRStatus(ctx context.Context, prID, newPrStatus string) error
	MergePr(ctx context.Context, prID string) (*entity.PullRequest, error)
	ReassignPrReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (entity.PullRequest, error)
	GetTeamByUserID(ctx context.Context, userID string) (string, error)
	CheckTeam(ctx context.Context, teamName string) (bool, error)
	CheckUser(ctx context.Context, userID string) (bool, error)
	CheckPR(ctx context.Context, prID string) (bool, error)
}

// UseCase - бизнес логика
type UseCase struct {
	repo RepositoryProvider
}

// New - конструктор бизнес логики
func New(repo RepositoryProvider) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// CreateTeam - создание команды
func (uc *UseCase) CreateTeam(ctx context.Context, team entity.Team) (*entity.Team, error) {
	if len(team.Members) == 0 {
		return nil, fmt.Errorf("members is empty")
	}

	if team.TeamName == "" {
		return nil, fmt.Errorf("team name is empty")
	}

	// проверяем существование команды
	existTeam, err := uc.repo.CheckTeam(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}

	if existTeam {
		return nil, entity.ErrTeamExists
	}

	err = uc.repo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// GetTeam - получить название команды и участников
func (uc *UseCase) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	if teamName == "" {
		return nil, fmt.Errorf("team name is empty")
	}

	team, err := uc.repo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	if len(team.Members) == 0 {
		return nil, entity.ErrNotFound
	}

	return team, nil
}

// ChangeActivityUser - изменение активности пользователя
func (uc *UseCase) ChangeActivityUser(ctx context.Context, user entity.User) (*entity.User, error) {
	if user.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	existUser, err := uc.repo.CheckUser(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	if !existUser {
		return nil, entity.ErrNotFound
	}

	err = uc.repo.ChangeActivityUser(ctx, user.IsActive, user.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetReviewFromUser - получить pr для пользователя
func (uc *UseCase) GetReviewFromUser(ctx context.Context, userID string) ([]entity.PullRequestShort, error) {
	if userID == "" {
		return nil, fmt.Errorf("username is empty")
	}

	return uc.repo.GetReviewFromUser(ctx, userID)
}

// CreatePullRequest - создать pr
func (uc *UseCase) CreatePullRequest(ctx context.Context, pr entity.PullRequestShort) (*entity.PullRequest, error) {
	var fullPr entity.PullRequest

	if pr.PullRequestID == "" || pr.PullRequestName == "" {
		return nil, fmt.Errorf("pull request id or pull request name is empty")
	}

	// проверяем существование такого pr
	existPR, err := uc.repo.CheckPR(ctx, pr.PullRequestID)
	if err != nil {
		return nil, err
	}

	if existPR {
		return nil, entity.ErrPrExists
	}

	// проверяем существование пользователя
	existUser, err := uc.repo.CheckUser(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	if !existUser {
		return nil, entity.ErrNotFound
	}

	// смотрим в какой команде пользователь
	teamName, err := uc.repo.GetTeamByUserID(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers, err := uc.generateReviewers(ctx, teamName, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	// заполняем структуры pr-а
	fullPr.PullRequestID = pr.PullRequestID
	fullPr.PullRequestName = pr.PullRequestName
	fullPr.AuthorID = pr.AuthorID
	fullPr.Status = "OPEN"
	fullPr.AssignedReviewers = reviewers
	fullPr.CreatedAt = time.Now()

	err = uc.repo.CreatePullRequest(ctx, fullPr)
	if err != nil {
		return nil, err
	}

	return &fullPr, nil
}

// generateReviewers - генерация ревьюеров на pr
func (uc *UseCase) generateReviewers(ctx context.Context, teamName, authorID string) ([]string, error) {
	var candidates []string

	team, err := uc.repo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	for _, member := range team.Members {
		if member.IsActive && member.UserID != authorID {
			candidates = append(candidates, member.UserID)
		}
	}

	// если кандидатов нету возвращаем пустой список
	if len(candidates) == 0 {
		return candidates, nil
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > 2 {
		candidates = candidates[:2]
	}

	return candidates, nil
}

// MergePr - замержить pr
func (uc *UseCase) MergePr(ctx context.Context, prID string) (*entity.PullRequest, error) {
	if prID == "" {
		return nil, fmt.Errorf("prID is empty")
	}

	// проверяем есть ли такой pr
	existPR, err := uc.repo.CheckPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	if !existPR {
		return nil, entity.ErrNotFound
	}

	mergedPR, err := uc.repo.MergePr(ctx, prID)
	if err != nil {
		return nil, err
	}

	return mergedPR, nil
}

// ReassignPrReviewer - заменить ревьювера
func (uc *UseCase) ReassignPrReviewer(ctx context.Context, prID, oldReviewerID string) (*entity.PullRequest, string, error) {
	if prID == "" {
		return nil, "", fmt.Errorf("prID is empty")
	}
	if oldReviewerID == "" {
		return nil, "", fmt.Errorf("oldReviewerID is empty")
	}

	// Проверка, существует ли PR
	existPR, err := uc.repo.CheckPR(ctx, prID)
	if err != nil {
		return nil, "", err
	}
	if !existPR {
		return nil, "", entity.ErrNotFound
	}

	// Получаем PR для проверки статуса и ревьюверов
	checkPr, err := uc.repo.GetPR(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	// Проверка на merge pr-а
	if checkPr.Status != "OPEN" {
		return nil, "", entity.ErrPrMerged
	}

	// Проверка, что oldReviewerID действительно назначен на этот PR
	found := false
	for _, reviewerID := range checkPr.AssignedReviewers {
		if reviewerID == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", fmt.Errorf("reviewer is not assigned to this PR")
	}

	// Проверка, существует ли старый ревьювер
	existUser, err := uc.repo.CheckUser(ctx, oldReviewerID)
	if err != nil {
		return nil, "", err
	}
	if !existUser {
		return nil, "", entity.ErrNotFound
	}

	// Берем команду старого ревьювера
	teamName, err := uc.repo.GetTeamByUserID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", err
	}

	// Исключаем: старого ревьювера, автора PR и всех уже назначенных ревьюверов
	excludeIDs := []string{oldReviewerID, checkPr.AuthorID}
	// Исключаем всех текущих ревьюверов (кроме старого, он уже в списке)
	for _, reviewerID := range checkPr.AssignedReviewers {
		if reviewerID != oldReviewerID {
			excludeIDs = append(excludeIDs, reviewerID)
		}
	}

	// Генерируем нового ревьювера
	newReviewerID, err := uc.selectNewReviewer(ctx, teamName, excludeIDs...)
	if err != nil {
		return nil, "", err
	}

	// Обновляем PR в репозитории
	pr, err := uc.repo.ReassignPrReviewer(ctx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return nil, "", err
	}

	return &pr, newReviewerID, nil
}

// selectNewReviewer - выбрать нового ревьера, исключая указанные ID
func (uc *UseCase) selectNewReviewer(ctx context.Context, teamName string, excludeIDs ...string) (string, error) {
	team, err := uc.GetTeam(ctx, teamName)
	if err != nil {
		return "", err
	}

	excludeMap := make(map[string]struct{})
	for _, id := range excludeIDs {
		excludeMap[id] = struct{}{}
	}

	var candidates []string
	for _, m := range team.Members {
		if m.IsActive {
			if _, excluded := excludeMap[m.UserID]; !excluded {
				candidates = append(candidates, m.UserID)
			}
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("no active replacement candidate in team")
	}

	return candidates[rand.Intn(len(candidates))], nil
}
