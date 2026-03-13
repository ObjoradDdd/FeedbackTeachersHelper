package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type StudentHandler struct {
	StudentService *services.StudentService
}

func NewStudentHandler(studentService *services.StudentService) *StudentHandler {
	return &StudentHandler{
		StudentService: studentService,
	}
}

// CreateStudent godoc
// @Summary Create student
// @Description Creates student in a group
// @Tags students
// @Accept json
// @Produce json
// @Security UserID
// @Param input body dto.CreateStudentRequest true "Student payload"
// @Success 201 {object} dto.CreateStudentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /students [post]
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := GetUserID(w, r)
	if err != nil {
		return
	}

	var req dto.CreateStudentRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	studentId, err := h.StudentService.CreateStudent(req.Name, req.GroupId, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateStudentResponse{
		Id: studentId,
	})
}

// GetStudentsGroup godoc
// @Summary List group students
// @Description Returns students for group by group id
// @Tags students
// @Produce json
// @Security UserID
// @Param groupId path int true "Group ID"
// @Success 200 {object} dto.GetStudentsGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /students/{groupId} [get]
func (h *StudentHandler) GetStudentsGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := GetUserID(w, r)
	if err != nil {
		return
	}

	groupIdStr := r.PathValue("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Invalid group ID"})
		return
	}

	students, err := h.StudentService.GetGroupStudents(groupId, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetStudentsGroupResponse{
		Students: func() []dto.StudentDto {
			dtoStudents := make([]dto.StudentDto, len(students))
			for i, s := range students {
				dtoStudents[i] = *s.ToDto()
			}
			return dtoStudents
		}(),
	})
}

// UpdateStudent godoc
// @Summary Update student
// @Description Updates student by student id
// @Tags students
// @Accept json
// @Produce json
// @Security UserID
// @Param id path int true "Student ID"
// @Param input body dto.UpdateStudentRequest true "Student payload"
// @Success 200 {object} dto.UpdateStudentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /students/{id} [put]
func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := GetUserID(w, r)
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

	var req dto.UpdateStudentRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.StudentService.UpdateStudent(id, req.Name, req.GroupId, userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.UpdateStudentResponse{
		Message: "Student updated successfully",
	})
}

// DeleteStudent godoc
// @Summary Delete student
// @Description Deletes student by student id
// @Tags students
// @Produce json
// @Security UserID
// @Param id path int true "Student ID"
// @Success 200 {object} dto.DeleteStudentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /students/{id} [delete]
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := GetUserID(w, r)
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

	if err := h.StudentService.DeleteStudent(id, userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteStudentResponse{
		Message: "Student deleted successfully",
	})
}
