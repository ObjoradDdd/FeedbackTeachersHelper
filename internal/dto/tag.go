package dto

type TagDto struct {
	Name    string `json:"name"`
	Meaning string `json:"meaning"`
	Id      int    `json:"id"`
}

type GetTeacherTagsRequest struct {
	Tags []TagDto `json:"tags"`
}

type CreateTagRequest struct {
	Name    string `json:"name"`
	Meaning string `json:"meaning"`
}

type CreateTagResponse struct {
	Id int `json:"id"`
}

type DeleteTagRequest struct {
	Id int `json:"id"`
}

type DeleteTagResponse struct {
	Id int `json:"id"`
}

type UpdateTagRequest struct {
	Name    string `json:"name"`
	Meaning string `json:"meaninig"`
}

type UpdateTagResponse struct {
	Id int `json:"id"`
}
