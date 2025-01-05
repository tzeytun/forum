package structs

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userid"`
	UserName  string `json:"username"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	LikeCount int    `json:"likecount"`
	PhotoPath string `json:"photopath"`
}

type Comment struct {
	ID        int    `json:"id"`
	PostId    int    `json:"postid"`
	UserId    int    `json:"userid"`
	UserName  string `json:"username"`
	Comment   string `json:"comment"`
	LikeCount int    `json:"likecount"`
}

type PostWithComments struct {
	Post     Post
	Comments []Comment
}
