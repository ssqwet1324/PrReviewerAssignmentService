package usecase

import (
	"context"
	"pr_reviewer_service/internal/entity"
	pkgmetrics "pr_reviewer_service/pkg/prometheus"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const nameTracer = "usecase"

// UseCaseObs - структура для метрик usecase
type UseCaseObs struct {
	UseCase
	metrics *pkgmetrics.Metrics
}

// NewObs - конструктор метрик
func NewObs(uc UseCase) *UseCaseObs {
	return &UseCaseObs{
		UseCase: uc,
		metrics: pkgmetrics.NewMetrics("usecase"),
	}
}

// CreateTeam - метрики
func (uc *UseCaseObs) CreateTeam(ctx context.Context, team entity.Team) (*entity.Team, error) {
	const methodName = "create_team"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp, err := uc.UseCase.CreateTeam(ctx, team)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.CreateTeam")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp, err
}

// GetTeam - метрики
func (uc *UseCaseObs) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	const methodName = "get_team"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp, err := uc.UseCase.GetTeam(ctx, teamName)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.GetTeam")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp, err
}

// ChangeActivityUser - метрики
func (uc *UseCaseObs) ChangeActivityUser(ctx context.Context, user entity.User) (*entity.User, error) {
	const methodName = "change_activity_user"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp, err := uc.UseCase.ChangeActivityUser(ctx, user)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.ChangeActivityUser")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp, err
}

// GetReviewFromUser  - метрики
func (uc *UseCaseObs) GetReviewFromUser(ctx context.Context, userID string) ([]entity.PullRequestShort, error) {
	const methodName = "get_review_from_user"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp, err := uc.UseCase.GetReviewFromUser(ctx, userID)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.GetReviewFromUser")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp, err
}

// CreatePullRequest - метрики
func (uc *UseCaseObs) CreatePullRequest(ctx context.Context, pr entity.PullRequestShort) (*entity.PullRequest, error) {
	const methodName = "create_pull_request"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp, err := uc.UseCase.CreatePullRequest(ctx, pr)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.CreatePullRequest")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp, err
}

// MergePr - метрики
func (uc *UseCaseObs) MergePr(ctx context.Context, prID string) (*entity.PullRequest, error) {
	const methodName = "merge_pr"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp, err := uc.UseCase.MergePr(ctx, prID)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.MergePr")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp, err
}

// ReassignPrReviewer - метрики
func (uc *UseCaseObs) ReassignPrReviewer(ctx context.Context, prID, oldReviewerID string) (*entity.PullRequest, string, error) {
	const methodName = "reassign_pr_reviewer"

	tracer := otel.Tracer(nameTracer)
	_, span := tracer.Start(ctx, methodName)
	defer span.End()

	startTime := time.Now()

	resp1, resp2, err := uc.UseCase.ReassignPrReviewer(ctx, prID, oldReviewerID)
	if err != nil {
		uc.metrics.HitError(methodName)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to uc.UseCase.ReassignPrReviewer")
	} else {
		uc.metrics.HitSuccess(methodName)
	}

	uc.metrics.HitDuration(methodName, time.Since(startTime).Seconds())

	return resp1, resp2, err
}
