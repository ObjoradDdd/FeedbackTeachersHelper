package models

type Teacher struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Student struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Group struct {
	Students  []Student `json:"students"`
	TeacherId int       `json:"teacher_id"`
	Name      string    `json:"name"`
	ID        int       `json:"id"`
}

type Tag struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	IsBad bool   `json:"is_bad"`
}

type StudentFeedback struct {
	Student string `json:"student_name"`
	TagIDs  []int  `json:"tag_ids"`
}

type FeedbackRequest struct {
	Theme     string            `json:"theme"`
	WhatWeDid string            `json:"what_we_did"`
	Feedback  []StudentFeedback `json:"feedback"`
}

type StudentFeedbackGenerated struct {
	Student string `json:"student_name"`
	Text    string `json:"text"`
}

type FeedbackResponse struct {
	Feedback []StudentFeedbackGenerated `json:"feedback"`
}
