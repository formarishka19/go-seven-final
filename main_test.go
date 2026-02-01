package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 100},
	}
	for city, list := range cafeList {
		for _, value := range requests {
			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&count=%d", city, value.count), nil)
			handler.ServeHTTP(response, req)

			require.Equal(t, http.StatusOK, response.Code)
			result := strings.TrimSpace(response.Body.String())

			if len(result) != 0 {
				reultCafeList := strings.Split(result, ",")
				resultCafeCount := len(reultCafeList)
				if value.want != 100 {
					assert.Equal(t, value.want, resultCafeCount)
				} else {
					assert.Equal(t, len(list), resultCafeCount)

				}

			} else {
				assert.Equal(t, value.want, 0)
			}
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, value := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", value.search), nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
		result := strings.TrimSpace(response.Body.String())
		count := 0
		if len(result) == 0 {
			assert.Equal(t, value.wantCount, count)
		} else {
			resultCafeList := strings.Split(result, ",")
			for _, s := range resultCafeList {
				if strings.Contains(strings.ToLower(s), strings.ToLower(value.search)) {
					count++
				}
			}
			assert.Equal(t, value.wantCount, count)
		}
	}
}
