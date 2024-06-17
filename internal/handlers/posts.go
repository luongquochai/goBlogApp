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

	var post models.CreatePostParams

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.AuthorID = sql.NullInt32{
		Int32: user_id.(int32),
		Valid: true,
	}

	// log.Printf("Post data: %+v\n", post)

	result, err := models.New(p.db).CreatePost(c, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)

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

func (p *PostHandler) ListPosts(c *gin.Context) {
	posts, err := models.New(p.db).ListPosts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (p *PostHandler) UpdatePostByAuthor(c *gin.Context) {
	var params models.UpdatePostByAuthorParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	params.ID = int32(id)
	params.AuthorID.Int32 = int32(c.GetInt("author_id"))

	err := models.New(p.db).UpdatePostByAuthor(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"updated": params.Content,
	})

}

func (p *PostHandler) DeletePost(c *gin.Context) {
	var params models.DeletePostByAuthorParams
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	authorID, err := strconv.Atoi(c.PostForm("author_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	params.ID = int32(id)
	params.AuthorID.Int32 = int32(authorID)

	err = models.New(p.db).DeletePostByAuthor(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"deleted": params.ID,
	})
}
