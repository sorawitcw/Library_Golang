package handler

import (
    "library/core/entity"
    "library/core/ports"
    "net/http"

    "github.com/gin-gonic/gin"
)

type MemberHandler struct {
    service ports.MemberService
}

func NewMemberHandler(service ports.MemberService) *MemberHandler {
    return &MemberHandler{service: service}
}

func (h *MemberHandler) Register(c *gin.Context) {
    var body entity.Member
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := h.service.Register(body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, result)
}

func (h *MemberHandler) Login(c *gin.Context) {
    var body struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.service.Login(body.Username, body.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *MemberHandler) GetProfile(c *gin.Context) {
    // ดึง memberId จาก JWT claim ที่ middleware เซ็ตไว้
    memberId, _ := c.Get("memberId")

    profile, err := h.service.GetProfile(memberId.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, profile)
}

func (h *MemberHandler) GetAllMembers(c *gin.Context) {
    role, _ := c.Get("role")

    members, err := h.service.GetAllMembers(role.(string))
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, members)
}