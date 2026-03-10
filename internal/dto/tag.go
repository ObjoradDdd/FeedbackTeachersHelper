package dto

type TagDto struct {
	Name    string
	Meaning string
	Id      int
}

type GetTeacherTagsRequest struct {
	Tags []TagDto
}

type CreateTagRequest struct {
	Name    string
	Meaning string
}

type CreateTagResponse struct {
	Id int
}
