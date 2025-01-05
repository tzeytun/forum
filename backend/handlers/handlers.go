package handlers

import (
	"net/http"

	createcomment "forum/backend/controllers/create/createComment"
	createpost "forum/backend/controllers/create/createPost"
	deleteaccount "forum/backend/controllers/delete/deleteAccount"
	deletecomment "forum/backend/controllers/delete/deleteComment"
	deletepost "forum/backend/controllers/delete/deletePost"
	getallposts "forum/backend/controllers/get/getAllPosts"
	getmycomments "forum/backend/controllers/get/getMyComments"
	getmyposts "forum/backend/controllers/get/getMyPosts"
	getmyvotedposts "forum/backend/controllers/get/getMyVotedPosts"
	getpostandcomments "forum/backend/controllers/get/getPostAndComments"
	getsearchedposts "forum/backend/controllers/get/getSearchedPosts"
	"forum/backend/controllers/login"
	"forum/backend/controllers/logout"
	"forum/backend/controllers/register"
	downvote "forum/backend/controllers/votes/downVote"
	upvote "forum/backend/controllers/votes/upVote"
	createpostpage "forum/frontend/pages/createPostPage"
	deleteaccountpage "forum/frontend/pages/deleteAccountPage"
	loginpage "forum/frontend/pages/loginPage"
	mainpage "forum/frontend/pages/mainPage"
	postpage "forum/frontend/pages/postPage"
	mycommentspage "forum/frontend/pages/profile/myCommentsPage"
	mypostspage "forum/frontend/pages/profile/myPostsPage"
	myvotedpostspage "forum/frontend/pages/profile/myVotedPostsPage"
	registerpage "forum/frontend/pages/registerPage"
	searchedpostspage "forum/frontend/pages/searchedPostsPage"
)

func ImportHandlers() {
	// API
	http.HandleFunc("/api/register", register.Register)
	http.HandleFunc("/api/login", login.Login)
	http.HandleFunc("/api/logout", logout.Logout)
	http.HandleFunc("/api/createpost", createpost.CreatePost)
	http.HandleFunc("/api/createcomment", createcomment.CreateComment)
	http.HandleFunc("/api/deleteaccount", deleteaccount.DeleteAccount)
	http.HandleFunc("/api/deletepost", deletepost.DeletePost)
	http.HandleFunc("/api/deletecomment", deletecomment.DeleteComment)
	http.HandleFunc("/api/upvote", upvote.UpVote)
	http.HandleFunc("/api/downvote", downvote.DownVote)
	http.HandleFunc("/api/allposts", getallposts.GetAllPosts)
	http.HandleFunc("/api/postandcomments", getpostandcomments.GetPostAndComments)
	http.HandleFunc("/api/myposts", getmyposts.GetMyPosts)
	http.HandleFunc("/api/mycomments", getmycomments.GetMyComments)
	http.HandleFunc("/api/myvotedposts", getmyvotedposts.GetMyVotedPosts)
	http.HandleFunc("/api/searchedposts", getsearchedposts.GetSearchedPosts)

	// Front-end
	http.HandleFunc("/", mainpage.MainPage)
	http.HandleFunc("/register", registerpage.RegisterPage)
	http.HandleFunc("/login", loginpage.LoginPage)
	http.HandleFunc("/createpost", createpostpage.CreatePostPage)
	http.HandleFunc("/post", postpage.PostPage)
	http.HandleFunc("/createcomment", postpage.PostPageCreateComment)
	http.HandleFunc("/upvote", postpage.PostPageUpVote)
	http.HandleFunc("/downvote", postpage.PostPageDownVote)
	http.HandleFunc("/deleteaccount", deleteaccountpage.DeleteAccountPage)
	http.HandleFunc("/myposts", mypostspage.MyPostsPage)
	http.HandleFunc("/deletepost", mypostspage.DeleteMyPost)
	http.HandleFunc("/mycomments", mycommentspage.MyCommentsPage)
	http.HandleFunc("/deletecomment", mycommentspage.DeleteMyComment)
	http.HandleFunc("/myvotedposts", myvotedpostspage.MyVotedPostsPage)
	http.HandleFunc("/search", searchedpostspage.SearchedPostsPage)
	http.HandleFunc("/login/google", login.HandleGoogleLogin)
	http.HandleFunc("/callback/google", login.HandleGoogleCallback)
	http.HandleFunc("/login/github", login.HandleGitHubLogin)
	http.HandleFunc("/callback/github", login.HandleGitHubCallback)
	http.HandleFunc("/login/facebook", login.HandleFacebookLogin)
	http.HandleFunc("/callback/facebook", login.HandleFacebookCallback)

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))
}
