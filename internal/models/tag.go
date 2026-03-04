package models

type Tag struct {
	Name      string `json:"name"`
	Meaning   string `json:"meaning"`
	ID        int    `json:"id"`
	TeacherID int    `json:"teacher_id"`
	IsBad     bool   `json:"is_bad"`
}
