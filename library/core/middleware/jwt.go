package middleware

import (
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
)

var jwtSecret = []byte("library-secret-key")

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ดึง token จาก Header: Authorization: Bearer <token>
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "กรุณา login ก่อน",
            })
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "รูปแบบ token ไม่ถูกต้อง",
            })
            return
        }

        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "token ไม่ถูกต้องหรือหมดอายุ",
            })
            return
        }

        // เก็บ claim ไว้ใน context เพื่อให้ handler ดึงไปใช้ต่อ
        claims := token.Claims.(jwt.MapClaims)
        c.Set("memberId", claims["memberId"])
        c.Set("role", claims["role"])

        c.Next()
    }
}

// RequireRole — ใช้ต่อจาก JWTMiddleware เพื่อตรวจ role
func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, _ := c.Get("role")
        if userRole != role {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "ไม่มีสิทธิ์เข้าถึง",
            })
            return
        }
        c.Next()
    }
}