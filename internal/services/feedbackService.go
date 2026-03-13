package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type FeedbackStorage interface {
	GetApiKey(userID int) (string, error)
	GetGroupStudents(groupID int, userID int) ([]models.Student, error)
	GetUserTags(userID int) ([]models.Tag, error)
}

type LlmClient interface {
	GenerateFeedback(ctx context.Context, prompt string, apiKey string) (string, error)
}

type FeedbackService struct {
	storage FeedbackStorage
	llm     LlmClient
}

func NewFeedbackService(storage FeedbackStorage, llm LlmClient) *FeedbackService {
	return &FeedbackService{storage: storage, llm: llm}
}

type StudentFeedbackInput struct {
	StudentId int
	Comment   string
	TagIds    []int
}

type GenerateFeedbackInput struct {
	GroupID           int
	LessonDescription string
	Activities        string
	Students          []StudentFeedbackInput
}

type StudentFeedback struct {
	StudentId int
	Name      string
	Comment   string
	Tags      []models.Tag
}

type GroupFeedback struct {
	UserID            int
	GroupID           int
	LessonDescription string
	Activities        string
	Students          []StudentFeedback
}

func (s *FeedbackService) GenerateFeedback(ctx context.Context, req *GenerateFeedbackInput, userID int) (models.GeneratedGroupFeedback, error) {
	hash, err := s.storage.GetApiKey(userID)

	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching API key: %w", err)
	}

	apiKey, err := Decrypt(hash)

	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching API key: %w", err)
	}

	groupStudents, err := s.storage.GetGroupStudents(req.GroupID, userID)
	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching students: %w", err)
	}

	userTags, err := s.storage.GetUserTags(userID)
	if err != nil {
		return models.GeneratedGroupFeedback{}, fmt.Errorf("error fetching tags: %w", err)
	}

	studentMap := make(map[int]string)
	for _, student := range groupStudents {
		studentMap[student.Id] = student.Name
	}

	tagMap := make(map[int]models.Tag)
	for _, tag := range userTags {
		tagMap[tag.Id] = tag
	}

	internalInput := &GroupFeedback{
		UserID:            userID,
		GroupID:           req.GroupID,
		LessonDescription: req.LessonDescription,
		Activities:        req.Activities,
		Students:          make([]StudentFeedback, 0, len(req.Students)),
	}

	for _, studentReq := range req.Students {
		studentName, exists := studentMap[studentReq.StudentId]
		if !exists {
			return models.GeneratedGroupFeedback{}, fmt.Errorf("student with Id %d not found in group %d", studentReq.StudentId, req.GroupID)
		}

		var studentTags []models.Tag
		for _, tagId := range studentReq.TagIds {
			if tag, tagExists := tagMap[tagId]; tagExists {
				studentTags = append(studentTags, tag)
			}
		}

		internalInput.Students = append(internalInput.Students, StudentFeedback{
			StudentId: studentReq.StudentId,
			Name:      studentName,
			Comment:   studentReq.Comment,
			Tags:      studentTags,
		})
	}

	prompt := generatePrompt(internalInput)

	feedbackText, err := s.llm.GenerateFeedback(ctx, prompt, apiKey)
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
	sb.WriteString("---\n")
	sb.WriteString("Имя третьего остальных учеников, которых я отправлю ::: Твой отзыв он них в 2-4 предложениях.\n")
	sb.WriteString("===\n")
	sb.WriteString("Общее описание урока в 3-5 предложений.\n")

	fmt.Println(sb.String())
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
		UserID:            input.UserID,
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

		var studentId int
		for _, s := range input.Students {
			if strings.EqualFold(strings.TrimSpace(s.Name), name) {
				studentId = s.StudentId
				break
			}
		}

		if studentId != 0 {
			result.Students = append(result.Students, models.GeneratedStudentFeedback{
				StudentId: studentId,
				Name:      name,
				Feedback:  feedbackText,
			})
		}
	}

	return result, nil
}
