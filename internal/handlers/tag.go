package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type TagHandler struct {
	TagService *services.TagService
}

func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		TagService: tagService,
	}
}

func (h *TagHandler) GetTeacherTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	tags, err := h.TagService.GetTeachersTags(teacherId)
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

	tagId, err := h.TagService.CreateTag(services.CreateTagInput{
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
