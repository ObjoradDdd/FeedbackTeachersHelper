package services

import (
	"fmt"
	"strings"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/utils"
)

type FeedbackStorage interface {
	GetApiKey(teacherId int) (string, error)
	GetGroupStudents(groupId int, teacherId int) ([]models.Student, error)
	GetTeachersTags(teacherId int) ([]models.Tag, error)
}

type LlmClient interface {
	GenerateFeedback(prompt string, apiKey string) (string, error)
}

type FeedbackService struct {
	storage FeedbackStorage
	llm     LlmClient
}

func NewFeedbackService(storage FeedbackStorage, llm LlmClient) *FeedbackService {
	return &FeedbackService{storage: storage, llm: llm}
}

type StudentFeedbackInput struct {
	StudentID int
	Comment   string
	TagIDs    []int
}

type GenerateFeedbackInput struct {
	TeacherID         int
	GroupID           int
	LessonDescription string
	Activities        string
	Students          []StudentFeedbackInput
}

type StudentFeedback struct {
	StudentID int
	Name      string
	Comment   string
	Tags      []models.Tag
}

type GroupFeedback struct {
	TeacherID         int
	GroupID           int
	LessonDescription string
	Activities        string
	Students          []StudentFeedback
}

func (s *FeedbackService) GenerateFeedback(req *GenerateFeedbackInput, teacherId int) (models.GeneratedGroupFeedback, error) {
	hash, err := s.storage.GetApiKey(teacherId)

	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching API key: %w", err)
	}

	apiKey, err := utils.Decrypt(hash)

	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching API key: %w", err)
	}

	groupStudents, err := s.storage.GetGroupStudents(req.GroupID, teacherId)
	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching students: %w", err)
	}

	teacherTags, err := s.storage.GetTeachersTags(teacherId)
	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching tags: %w", err)
	}

	studentMap := make(map[int]string)
	for _, student := range groupStudents {
		studentMap[student.ID] = student.Name
	}

	tagMap := make(map[int]models.Tag)
	for _, tag := range teacherTags {
		tagMap[tag.ID] = tag
	}

	internalInput := &GroupFeedback{
		TeacherID:         teacherId,
		GroupID:           req.GroupID,
		LessonDescription: req.LessonDescription,
		Activities:        req.Activities,
		Students:          make([]StudentFeedback, 0, len(req.Students)),
	}

	for _, studentReq := range req.Students {
		studentName, exists := studentMap[studentReq.StudentID]
		if !exists {
			return models.GeneratedGroupFeedback{}, fmt.Errorf("student with ID %d not found in group %d", studentReq.StudentID, req.GroupID)
		}

		var studentTags []models.Tag
		for _, tagID := range studentReq.TagIDs {
			if tag, tagExists := tagMap[tagID]; tagExists {
				studentTags = append(studentTags, tag)
			}
		}

		internalInput.Students = append(internalInput.Students, StudentFeedback{
			StudentID: studentReq.StudentID,
			Name:      studentName,
			Comment:   studentReq.Comment,
			Tags:      studentTags,
		})
	}

	prompt := generatePrompt(internalInput)

	feedbackText, err := s.llm.GenerateFeedback(prompt, apiKey)
	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error generating feedback: %w", err)
	}

	groupFeedbackResult, err := mapOutput(feedbackText, internalInput)
	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error mapping LLM output: %w", err)
	}

	return *groupFeedbackResult, nil
}

func generatePrompt(groupInput *GroupFeedback) string {
	var sb strings.Builder

	sb.WriteString("Ты — профессиональный и тактичный преподаватель. Помни, ты пишешь не ученикам, а родителям, не надо обращений, просто пиши что-то по типу этого: 'Иван сегодня молодец и тд' . Твоя задача — написать индивидуальные комментарии для родителей каждого ученика по итогам прошедшего занятия. В тексте не должно быть орфографических ошибок. Используй имена учеников, чтобы текст звучал мягче.\n\n")

	sb.WriteString(fmt.Sprintf("Тема урока: %s\n", groupInput.LessonDescription))

	if groupInput.Activities != "" {
		sb.WriteString(fmt.Sprintf("Что мы делали: %s\n", groupInput.Activities))
	}

	sb.WriteString("\nНиже представлен список учеников и теги, описывающие их поведение на уроке:\n")

	for _, student := range groupInput.Students {
		tagValues := make([]string, len(student.Tags))
		for i, tag := range student.Tags {
			tagValues[i] = tag.Meaning
		}
		tagsStr := strings.Join(tagValues, ", ")
		sb.WriteString(fmt.Sprintf("- %s (Теги: %s)\n", student.Name, tagsStr))

		if student.Comment != "" {
			sb.WriteString(fmt.Sprintf("  Комментарий учителя: %s\n", student.Comment))
		}
	}

	sb.WriteString("\nВЫДАЙ ОТВЕТ СТРОГО ПО ЭТОМУ ШАБЛОНУ (не добавляй лишнего текста до или после):\n")
	sb.WriteString("Имя первого ученика ::: Твой отзыв на 2-4 предложения.\n")
	sb.WriteString("---\n")
	sb.WriteString("Имя второго ученика ::: Твой отзыв на 2-4 предложения.\n")
	sb.WriteString("===\n")
	sb.WriteString("Общее описание урока в 3-5 предложений.\n")

	return sb.String()
}

func mapOutput(llmOutput string, input *GroupFeedback) (*models.GeneratedGroupFeedback, error) {
	parts := strings.Split(llmOutput, "===")
	if len(parts) < 2 {
		return nil, fmt.Errorf("не найден разделитель '==='. Ответ ИИ: %s", llmOutput)
	}

	studentsPart := parts[0]
	lessonDesc := strings.TrimSpace(parts[1])

	result := &models.GeneratedGroupFeedback{
		TeacherID:         input.TeacherID,
		GroupID:           input.GroupID,
		LessonDescription: lessonDesc,
		Students:          make([]models.GeneratedStudentFeedback, 0, len(input.Students)),
	}

	studentChunks := strings.Split(studentsPart, "---")

	for _, chunk := range studentChunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		nameAndFeedback := strings.SplitN(chunk, ":::", 2)
		if len(nameAndFeedback) < 2 {
			continue
		}

		name := strings.TrimSpace(nameAndFeedback[0])
		feedbackText := strings.TrimSpace(nameAndFeedback[1])

		name = strings.ReplaceAll(name, "*", "")
		name = strings.ReplaceAll(name, "#", "")

		var studentID int
		for _, s := range input.Students {
			if strings.EqualFold(strings.TrimSpace(s.Name), name) {
				studentID = s.StudentID
				break
			}
		}

		if studentID != 0 {
			result.Students = append(result.Students, models.GeneratedStudentFeedback{
				StudentID: studentID,
				Name:      name,
				Feedback:  feedbackText,
			})
		}
	}

	return result, nil
}
