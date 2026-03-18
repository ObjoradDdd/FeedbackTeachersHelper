package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/go-playground/validator/v10"
)

type groupService interface {
	CreateGroup(ctx context.Context, name string, userID int) (int, error)
	GetUserGroups(ctx context.Context, userID int) ([]models.Group, error)
	UpdateGroup(ctx context.Context, id int, name string, userID int) error
	DeleteGroup(ctx context.Context, id int, userID int) error
}

type GroupHandler struct {
	groupService groupService
	validator    *validator.Validate
}

func NewGroupHandler(groupService groupService, validator *validator.Validate) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		validator:    validator,
	}
}

// CreateGroup godoc
// @Summary Create group
// @Description Creates a group for current user
// @Tags groups
// @Accept json
// @Produce json
// @Security UserID
// @Param input body dto.CreateGroupRequest true "Group payload"
// @Success 200 {object} dto.CreateGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	var req dto.CreateGroupRequest
	if err := decodeAndValidateRequest(w, r, &req, h.validator); err != nil {
		return
	}

	groupId, err := h.groupService.CreateGroup(r.Context(), req.Name, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create group")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateGroupResponse{
		Id: groupId,
	})
}

// GetGroups godoc
// @Summary List groups
// @Description Returns all groups for current user
// @Tags groups
// @Produce json
// @Security UserID
// @Success 200 {object} dto.GetGroupsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups [get]
func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	groups, err := h.groupService.GetUserGroups(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get groups")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.GetGroupsResponse{
		Groups: func() []dto.GroupDto {
			result := make([]dto.GroupDto, len(groups))
			for i, group := range groups {
				result[i] = *group.ToDto()
			}
			return result
		}(),
	})
}

// UpdateGroup godoc
// @Summary Update group
// @Description Updates group name by group id
// @Tags groups
// @Accept json
// @Produce json
// @Security UserID
// @Param id path int true "Group ID"
// @Param input body dto.UpdateGroupRequest true "Group payload"
// @Success 200 {object} dto.UpdateGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups/{id} [put]
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdateGroupRequest
	if err := decodeAndValidateRequest(w, r, &req, h.validator); err != nil {
		return
	}

	err = h.groupService.UpdateGroup(r.Context(), id, req.Name, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update group")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.UpdateGroupResponse{
		Message: "Group updated successfully",
	})
}

// DeleteGroup godoc
// @Summary Delete group
// @Description Deletes group by group id
// @Tags groups
// @Produce json
// @Security UserID
// @Param id path int true "Group ID"
// @Success 200 {object} dto.DeleteGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if id < 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.groupService.DeleteGroup(r.Context(), id, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete group")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.DeleteGroupResponse{
		Message: "Group deleted successfully ",
	})
}
