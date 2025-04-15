package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoginAs(t *testing.T, client *http.Client, username, password string) APIResponse {
	return DoPost(t, client, "/auth", fmt.Sprintf("username=%s&password=%s", username, password))
}

func TestLogin(t *testing.T) {
	t.Run("Correct credentials", func(t *testing.T) {
		client := InitClient(t)
		r := LoginAs(t, client, "Smaug", "10082005")
		assert.Equal(t, "success", r.Status)
		assert.Equal(t, 200, r.StatusCode)
		assert.IsType(t, map[string]any{}, r.Result)
		data, ok := r.Result.(map[string]any)
		if !ok {
			t.Fatalf("The 'result' field is not of type JSON (map[string]any)")
		}
		// Check user data
		assert.Contains(t, data, "id")
		assert.IsType(t, float64(0), data["id"]) // JSON numbers → float64
		assert.Contains(t, data, "username")
		assert.IsType(t, "", data["username"])
		assert.Contains(t, data, "firstname")
		assert.IsType(t, "", data["firstname"])
		assert.Contains(t, data, "lastname")
		assert.IsType(t, "", data["lastname"])
		assert.Contains(t, data, "role")
		assert.IsType(t, float64(0), data["role"])
		assert.Contains(t, data, "avatar") // Can be nil
		assert.Contains(t, data, "email")
		assert.IsType(t, "", data["email"])
		assert.Contains(t, data, "created_timestamp")
		assert.IsType(t, float64(0), data["created_timestamp"])
		assert.Contains(t, data, "updated_timestamp")
		assert.IsType(t, float64(0), data["updated_timestamp"])
		// Check if token is present in cookies
		url, err := url.Parse(BaseURL)
		if err != nil {
			t.Fatalf("Failed to parse base URL: %v", err)
		}
		cookies := client.Jar.Cookies(url)
		assert.Len(t, cookies, 2)
		assert.Equal(t, "Authorization", cookies[0].Name)
		assert.Equal(t, "RefreshToken", cookies[1].Name)

	})
	t.Run("Incorrect credentials", func(t *testing.T) {
		client := InitClient(t)
		r := LoginAs(t, client, "Smaug", "wrongpassword")
		assert.Equal(t, "error", r.Status)
		assert.Equal(t, 401, r.StatusCode)
		assert.Equal(t, "Invalid credentials", r.Message)
	})
	t.Run("Missing username", func(t *testing.T) {
		client := InitClient(t)
		r := LoginAs(t, client, "", "password")
		assert.Equal(t, "error", r.Status)
		assert.Equal(t, 400, r.StatusCode)
		assert.Equal(t, "Username is required.", r.Message)
	})
	t.Run("Missing password", func(t *testing.T) {
		client := InitClient(t)
		r := LoginAs(t, client, "username", "")
		assert.Equal(t, "error", r.Status)
		assert.Equal(t, 400, r.StatusCode)
		assert.Equal(t, "Password is required.", r.Message)
	})
	t.Run("Missing credentials", func(t *testing.T) {
		client := InitClient(t)
		r := LoginAs(t, client, "", "")
		assert.Equal(t, "error", r.Status)
		assert.Equal(t, 400, r.StatusCode)
		assert.Equal(t, "Username is required. Password is required.", r.Message)
	})

}
