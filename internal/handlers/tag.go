package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// GetTeacherTags godoc
// @Summary List tags
// @Description Returns all tags for current teacher
// @Tags tags
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.GetTeacherTagsRequest
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag [get]
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

// CreateTag godoc
// @Summary Create tag
// @Description Creates tag for current teacher
// @Tags tags
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body dto.CreateTagRequest true "Tag payload"
// @Success 200 {object} dto.CreateTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag [post]
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

// DeleteTag godoc
// @Summary Delete tag
// @Description Deletes tag by tag id
// @Tags tags
// @Produce json
// @Security Bearer
// @Param id path int true "Tag ID"
// @Success 200 {object} dto.DeleteTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag/{id} [delete]
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
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

	err = h.tagService.DeleteTag(id, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteTagResponse{
		Id: id,
	})
}

// UpdateTag godoc
// @Summary Update tag
// @Description Updates tag by tag id
// @Tags tags
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Tag ID"
// @Param input body dto.UpdateTagRequest true "Tag payload"
// @Success 200 {object} dto.UpdateTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag/{id} [put]
func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateTagRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	err = h.tagService.UpdateTag(services.UpdateTagInput{
		Id:      id,
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
		Id: id,
	})
}
