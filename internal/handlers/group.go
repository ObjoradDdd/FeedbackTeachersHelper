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
