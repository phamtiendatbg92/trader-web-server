package dtos

type TutorialResponse struct {
	Meta Meta         `json:"meta"`
	Data TutorialJson `json:"data"`
}
type ListTutorialResponse struct {
	Meta *Meta         `json:"meta"`
	Data *TutorialJson `json:"data"`
}

type TutorialJson struct {
	Id          uint     `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	Url         string   `json:"url"`
	AuthorID    uint     `json:"authorId"`
	Description string   `json:"description"`
}
type Comment struct {
	UserId  uint   `json:"userId"`
	Comment string `json:"comment"`
}
type CommentItemRes struct {
	Id          uint      `json:"id"`
	CommentItem Comment   `json:"comment"`
	Replies     []Comment `json:"replies"`
}

type Comment1Req struct {
	PostId  uint   `json:"postId"`
	Comment string `json:"comment"`
}
type Comment2Req struct {
	ParentId uint   `json:"parentId"`
	Comment  string `json:"comment"`
}

type DeleteCommentReq struct {
	ID    uint `json:"id"`
	Level byte `json:"level"`
}
