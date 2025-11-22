package handler

import (
	"errors"
	"net/http"
	"pr_reviewer_service/internal/entity"
	"pr_reviewer_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Handler - ручки
type Handler struct {
	uc *usecase.UseCase
}

// New - конструктор handler
func New(uc *usecase.UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// CreateTeam - создать команду
func (h *Handler) CreateTeam(ctx *gin.Context) {
	var team entity.Team

	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Error: entity.ErrorDetail{
				Code:    "400",
				Message: err.Error(),
			},
		})
		return
	}

	newTeam, err := h.uc.CreateTeam(ctx, team)
	if err != nil {
		if errors.Is(err, entity.ErrTeamExists) {
			ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "TEAM_EXISTS",
					Message: "team_name already exists",
				},
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"team": newTeam,
	})
}

// GetTeam - получить команду
func (h *Handler) GetTeam(ctx *gin.Context) {
	teamName := ctx.Query("team_name")

	team, err := h.uc.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"team": team,
	})
}

// SetIsActive - изменить активность
func (h *Handler) SetIsActive(ctx *gin.Context) {
	var user entity.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Error: entity.ErrorDetail{
				Code:    "400",
				Message: err.Error(),
			},
		})
		return
	}

	data, err := h.uc.ChangeActivityUser(ctx, user)
	if err != nil {
		if data == nil {
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": data.UserID, "is_active": data.IsActive})
}

// GetReview - получить pr-ы где пользователь reviewer
func (h *Handler) GetReview(ctx *gin.Context) {
	userID := ctx.Query("user_id")

	pr, err := h.uc.GetReviewFromUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Error: entity.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID, "pull_requests": pr})
}

// PullRequestCreate - создание pr-а
func (h *Handler) PullRequestCreate(ctx *gin.Context) {
	var pr entity.PullRequestShort

	if err := ctx.ShouldBindJSON(&pr); err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Error: entity.ErrorDetail{
				Code:    "400",
				Message: err.Error(),
			},
		})
		return
	}

	fullPr, err := h.uc.CreatePullRequest(ctx, pr)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		} else if errors.Is(err, entity.ErrPrExists) {
			ctx.JSON(http.StatusConflict, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
				},
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"pr": fullPr})
}

// MergePR - замержить pr
func (h *Handler) MergePR(ctx *gin.Context) {
	var pr entity.PullRequestShort

	if err := ctx.ShouldBindJSON(&pr); err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Error: entity.ErrorDetail{
				Code:    "400",
				Message: err.Error(),
			},
		})
		return
	}

	mergedPr, err := h.uc.MergePr(ctx, pr.PullRequestID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"pr": mergedPr})
}

// ReassignPrReviewer - Переназначить конкретного ревьювера на другого из его команды
func (h *Handler) ReassignPrReviewer(ctx *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_reviewer_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Error: entity.ErrorDetail{
				Code:    "400",
				Message: err.Error(),
			},
		})
		return
	}

	pr, newReviewerID, err := h.uc.ReassignPrReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
		} else if errors.Is(err, entity.ErrPrMerged) {
			ctx.JSON(http.StatusConflict, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR",
				},
			})

		} else {
			ctx.JSON(http.StatusInternalServerError, entity.ErrorResponse{
				Error: entity.ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"pr": pr, "replaced_by": newReviewerID})
}
