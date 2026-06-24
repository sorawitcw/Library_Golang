package handler

import (
    "library/core/entity"
    "library/core/ports"
    "net/http"

    "github.com/gin-gonic/gin"
)

type BookHandler struct {
    service ports.BookService
}

func NewBookHandler(service ports.BookService) *BookHandler {
    return &BookHandler{service: service}
}

func (h *BookHandler) AddBook(c *gin.Context) {
    var body entity.Book
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := h.service.AddBook(body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, result)
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
    books, err := h.service.GetAllBooks()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBookById(c *gin.Context) {
    id := c.Param("id")
    book, err := h.service.GetBookById(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
    id := c.Param("id")
    var body entity.Book
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := h.service.UpdateBook(id, body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
    id := c.Param("id")
    if err := h.service.DeleteBook(id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "ลบหนังสือสำเร็จ"})
}