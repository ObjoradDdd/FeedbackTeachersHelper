package dto

type TagDto struct {
	Name    string `json:"name"`
	Meaning string `json:"meaning"`
	Id      int    `json:"id"`
}

type GetUserTagsResponse struct {
	Tags []TagDto `json:"tags"`
}

type CreateTagRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=64"`
	Meaning string `json:"meaning" validate:"required,min=8,max=512"`
}

type CreateTagResponse struct {
	Id int `json:"id"`
}

type DeleteTagRequest struct {
	Id int `json:"id" validate:"required,gt=0"`
}

type DeleteTagResponse struct {
	Id int `json:"id"`
}

type UpdateTagRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=64"`
	Meaning string `json:"meaning" validate:"required,min=8,max=512"`
}

type UpdateTagResponse struct {
	Id int `json:"id"`
}
