package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
	"github.com/go-playground/validator/v10"
)

type tagService interface {
	GetUserTags(ctx context.Context, userID int) ([]models.Tag, error)
	CreateTag(ctx context.Context, input services.CreateTagInput, userID int) (int, error)
	DeleteTag(ctx context.Context, id int, userID int) error
	UpdateTag(ctx context.Context, input services.UpdateTagInput, userID int) error
}

type TagHandler struct {
	tagService tagService
	validator  *validator.Validate
}

func NewTagHandler(tagService tagService, validator *validator.Validate) *TagHandler {
	return &TagHandler{
		tagService: tagService,
		validator:  validator,
	}
}

// GetUserTags godoc
// @Summary List tags
// @Description Returns all tags for current user
// @Tags tags
// @Produce json
// @Security UserID
// @Success 200 {object} dto.GetUserTagsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag [get]
func (h *TagHandler) GetUserTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	tags, err := h.tagService.GetUserTags(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get tags")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.GetUserTagsResponse{
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
// @Description Creates tag for current user
// @Tags tags
// @Accept json
// @Produce json
// @Security UserID
// @Param input body dto.CreateTagRequest true "Tag payload"
// @Success 200 {object} dto.CreateTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag [post]
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	var req dto.CreateTagRequest
	if err := decodeAndValidateRequest(w, r, &req, h.validator); err != nil {
		return
	}

	tagId, err := h.tagService.CreateTag(r.Context(), services.CreateTagInput{
		Name:    req.Name,
		Meaning: req.Meaning,
	}, userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create tag")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.CreateTagResponse{
		Id: tagId,
	})
}

// DeleteTag godoc
// @Summary Delete tag
// @Description Deletes tag by tag id
// @Tags tags
// @Produce json
// @Security UserID
// @Param id path int true "Tag ID"
// @Success 200 {object} dto.DeleteTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag/{id} [delete]
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
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

	err = h.tagService.DeleteTag(r.Context(), id, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete tag")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.DeleteTagResponse{
		Id: id,
	})
}

// UpdateTag godoc
// @Summary Update tag
// @Description Updates tag by tag id
// @Tags tags
// @Accept json
// @Produce json
// @Security UserID
// @Param id path int true "Tag ID"
// @Param input body dto.UpdateTagRequest true "Tag payload"
// @Success 200 {object} dto.UpdateTagResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tag/{id} [put]
func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateTagRequest
	if err := decodeAndValidateRequest(w, r, &req, h.validator); err != nil {
		return
	}

	err = h.tagService.UpdateTag(r.Context(), services.UpdateTagInput{
		Id:      id,
		Name:    req.Name,
		Meaning: req.Meaning,
	}, userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update tag")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.UpdateTagResponse{
		Id: id,
	})
}
