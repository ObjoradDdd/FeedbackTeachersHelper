package handlers

import (
	"encoding/json"
	"net/http"

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

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.CreateStudentRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	studentId, err := h.StudentService.CreateStudent(req.Name, req.GroupId, teacherId)
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

func (h *StudentHandler) GetStudentsGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.GetStudentsGroupRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	students, err := h.StudentService.GetGroupStudents(req.Id, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
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

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.UpdateStudentRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.StudentService.UpdateStudent(req.Id, req.Name, req.GroupId, teacherId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.UpdateStudentResponse{
		Message: "Student updated successfully",
	})
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.DeleteStudentRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	if err := h.StudentService.DeleteStudent(req.Id, teacherId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteStudentResponse{
		Message: "Student deleted successfully",
	})
}
