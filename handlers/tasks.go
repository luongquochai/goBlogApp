package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/luongquochai/goBlog/db"
	models "github.com/luongquochai/goBlog/db/models"
)

type Task struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

func GetTasks(c *gin.Context) {
	userID := getUserID(c)
	tasks, err := db.Queries.GetTaskByUserID(context.Background(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID := getUserID(c)
	createdTask, err := db.Queries.CreateTask(context.Background(), models.CreateTaskParams{
		Title:   task.Title,
		Content: task.Content,
		UserID:  int32(userID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

func UpdateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	userID := getUserID(c)

	err := db.Queries.UpdateTask(context.Background(), models.UpdateTaskParams{
		ID:      int32(id),
		Title:   task.Title,
		Content: task.Content,
		UserID:  int32(userID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := getUserID(c)

	err := db.Queries.DeleteTask(context.Background(), models.DeleteTaskParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted successfully"})
}

func getUserID(c *gin.Context) int {
	session, _ := store.Get(c.Request, "session")
	username := session.Values["username"].(string)
	user, _ := db.Queries.GetUserByUsername(context.Background(), username)
	return int(user.ID)
}
