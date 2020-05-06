package app

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"ticketing/auth/database"
)

func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := database.New()
	userRepo := database.NewUserRepository(db)

	tests := []struct {
		payload      string
		userRepo     *database.UserRepository
		expectedCode int
	}{
		{payload: "", expectedCode: http.StatusBadRequest},
		{payload: `{"email": "than2@amcil.com", "password": "hello"}`, expectedCode: http.StatusOK},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/", bytes.NewReader([]byte(tt.payload)))

		signUp(userRepo)(c)

		if w.Code != tt.expectedCode {
			t.Errorf("error status code, expect=%d but got=%d", tt.expectedCode, w.Code)
		}
	}
}

func TestSignOut(t *testing.T) {

}
