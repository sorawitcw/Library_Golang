package handler

import (
    "library/core/ports"
    "net/http"

    "github.com/gin-gonic/gin"
)

type BorrowHandler struct {
    service ports.BorrowService
}

func NewBorrowHandler(service ports.BorrowService) *BorrowHandler {
    return &BorrowHandler{service: service}
}

func (h *BorrowHandler) BorrowBook(c *gin.Context) {
    memberId, _ := c.Get("memberId")

    var body struct {
        BookId string `json:"book_id"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := h.service.BorrowBook(memberId.(string), body.BookId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, result)
}

func (h *BorrowHandler) ReturnBook(c *gin.Context) {
    memberId, _ := c.Get("memberId")
    recordId := c.Param("id")

    result, err := h.service.ReturnBook(recordId, memberId.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, result)
}

func (h *BorrowHandler) GetMyBorrows(c *gin.Context) {
    memberId, _ := c.Get("memberId")

    records, err := h.service.GetMyBorrows(memberId.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, records)
}