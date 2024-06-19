package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/luongquochai/goBlogApp/internal/models"
)

type PostHandler struct {
	db *sql.DB
}

func NewPostHandler(db *sql.DB) *PostHandler {
	return &PostHandler{db: db}
}

func (p *PostHandler) CreatePost(c *gin.Context) {
	session := sessions.Default(c)
	user_id := session.Get("user_id")
	// log.Printf("user_id: %v", user_id)

	if user_id == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse form data:
	title := c.PostForm("title")
	content := c.PostForm("content")
	authorID := sql.NullInt32{
		Int32: user_id.(int32),
		Valid: true,
	}

	if title == "" && content == "" {
		var post models.CreatePostParams
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		title = post.Title
		content = post.Content
	}

	// log.Printf("Post data: %+v\n", post)

	_, err := models.New(p.db).CreatePost(c.Request.Context(), models.CreatePostParams{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, result)
	c.Redirect(http.StatusSeeOther, "/home")

}

func (p *PostHandler) CancelHandler(c *gin.Context) {
	c.Redirect(http.StatusSeeOther, "/home")
}

func (p *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	post, err := models.New(p.db).GetPostByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (p *PostHandler) ListPostsHTML(c *gin.Context) {
	posts, err := models.New(p.db).ListPosts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assuming you have a 'posts.tmpl' file in your templates directory
	c.HTML(http.StatusOK, "home_auth.html", gin.H{
		"Posts": posts,
	})
}

func (p *PostHandler) ListPosts(c *gin.Context) {
	posts, err := models.New(p.db).ListPosts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)

	// c.HTML(http.StatusOK, "home_auth.html", gin.H{"posts": posts})
}

func (p *PostHandler) UpdatePostByAuthor(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var params models.UpdatePostByAuthorParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	params.ID = int32(id)
	params.AuthorID.Int32 = userID.(int32)
	params.AuthorID.Valid = true

	err := models.New(p.db).UpdatePostByAuthor(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"updated": params.Content,
	})
	c.Redirect(http.StatusSeeOther, "/")

}

func (p *PostHandler) DeletePost(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var params models.DeletePostByAuthorParams
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	params.ID = int32(id)
	params.AuthorID.Int32 = userID.(int32)
	params.AuthorID.Valid = true

	err = models.New(p.db).DeletePostByAuthor(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  "success",
	// 	"deleted": params.ID,
	// })
	c.Redirect(http.StatusSeeOther, "/home")
}
