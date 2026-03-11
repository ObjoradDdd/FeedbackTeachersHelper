package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup godoc
// @Summary Create group
// @Description Creates a group for current teacher
// @Tags groups
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body dto.CreateGroupRequest true "Group payload"
// @Success 200 {object} dto.CreateGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.CreateGroupRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	groupId, err := h.groupService.CreateGroup(req.Name, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateGroupResponse{
		Id: groupId,
	})
}

// GetGroups godoc
// @Summary List groups
// @Description Returns all groups for current teacher
// @Tags groups
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.GetGroupsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups [get]
func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	groups, err := h.groupService.GetTeachersGroups(teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetGroupsResponse{
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
// @Security Bearer
// @Param id path int true "Group ID"
// @Param input body dto.UpdateGroupRequest true "Group payload"
// @Success 200 {object} dto.UpdateGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups/{id} [put]
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Invalid ID"})
		return
	}

	var req dto.UpdateGroupRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	err = h.groupService.UpdateGroup(id, req.Name, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.UpdateGroupResponse{
		Message: "Group updated successfully",
	})
}

// DeleteGroup godoc
// @Summary Delete group
// @Description Deletes group by group id
// @Tags groups
// @Produce json
// @Security Bearer
// @Param id path int true "Group ID"
// @Success 200 {object} dto.DeleteGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Invalid ID"})
		return
	}

	err = h.groupService.DeleteGroup(id, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteGroupResponse{
		Message: "Group deleted successfully",
	})
}
