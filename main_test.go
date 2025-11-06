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
	city := "moscow"

	testCases := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList[city])},
	}

	for _, tc := range testCases {
		path := fmt.Sprintf("/cafe?city=%s&count=%d", city, tc.count)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := response.Body.String()
		if body == "" {
			assert.Equal(t, 0, tc.want)
			return
		}

		cafes := strings.Split(body, ",")
		assert.Equal(t, tc.want, len(cafes))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"

	testCases := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, tc := range testCases {
		path := fmt.Sprintf("/cafe?city=%s&search=%s", city, tc.search)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := response.Body.String()
		if body == "" {
			assert.Equal(t, 0, tc.wantCount)
			return
		}

		cafes := strings.Split(body, ",")
		assert.Equal(t, tc.wantCount, len(cafes))

		for _, cafe := range cafes {
			assert.True(t, strings.Contains(strings.ToLower(cafe), strings.ToLower(tc.search)))
		}
	}
}
