package main

import (
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

const erase_db = false

var todoCol *db.Col

type Todo struct {
	Id        string `json:"id"`
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed" binding:"required"`
}

func NewTodo(c *gin.Context) {
	var todo map[string]interface{}
	c.Bind(&todo)
	id, err := todoCol.Insert(todo)
	if err == nil {
		c.JSON(200, gin.H{"id": id})
	} else {
		log.Println(err)
		c.Writer.WriteHeader(500)
	}
}

func ListTodos(c *gin.Context) {
	todos := []Todo{}
	todoCol.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
		item := Todo{}
		json.Unmarshal(doc, &item)
		item.Id = strconv.Itoa(id)
		todos = append(todos, item)
		return true
	})
	c.JSON(200, todos)
}

func UpateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	var todo map[string]interface{}
	c.Bind(&todo)
	if err := todoCol.Update(id, todo); err == nil {
		c.Writer.WriteHeader(200)
	} else {
		log.Println(err)
		c.Writer.WriteHeader(404)
	}
}

func DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if err := todoCol.Delete(id); err == nil {
		c.Writer.WriteHeader(200)
	} else {
		log.Println(err)
		c.Writer.WriteHeader(404)
	}
}

func main() {
	myDBDir := "/tmp/TodoDatabase"

	if erase_db {
		os.RemoveAll(myDBDir)
		defer os.RemoveAll(myDBDir)
	}

	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		panic(err)
	}

	if erase_db {
		if err := myDB.Create("todo"); err != nil {
			panic(err)
		}
	}

	todoCol = myDB.Use("todo")

	r := gin.Default()
	r.Static("/todo", "public")

	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) { c.Writer.WriteHeader(200) })
		api.POST("/todos", NewTodo)
		api.GET("/todos", ListTodos)
		api.DELETE("/todos/:id", DeleteTodo)
		api.PUT("/todos/:id", UpateTodo)
	}

	r.Run(":8080")
}
