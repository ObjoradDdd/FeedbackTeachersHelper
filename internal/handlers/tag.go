package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type TagHandler struct {
	tagService *services.TagService
}

func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

func (h *TagHandler) GetTeacherTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	tags, err := h.tagService.GetTeachersTags(teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetTeacherTagsRequest{
		Tags: func(tags []models.Tag) []dto.TagDto {
			tagsDto := make([]dto.TagDto, len(tags))
			for i, tag := range tags {
				tagsDto[i] = tag.ToDto()
			}
			return tagsDto
		}(tags),
	})
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.CreateTagRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	tagId, err := h.tagService.CreateTag(services.CreateTagInput{
		Name:    req.Name,
		Meaning: req.Meaning,
	}, teacherId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateTagResponse{
		Id: tagId,
	})
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.DeleteTagRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	err = h.tagService.DeleteTag(req.Id, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteTagResponse{
		Id: req.Id,
	})
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.UpdateTagRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	err = h.tagService.UpdateTag(services.UpdateTagInput{
		Id:      req.Id,
		Name:    req.Name,
		Meaning: req.Meaning,
	}, teacherId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.UpdateTagResponse{
		Id: req.Id,
	})
}
