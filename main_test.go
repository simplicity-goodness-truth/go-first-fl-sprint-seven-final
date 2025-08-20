package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Structure for cafe count tests
type CafeCountReq struct {
	count int
	wait  int
}

type CafeSearchReq struct {
	search string
	wantCount   int
}

func TestCafeNotOk(t *testing.T) {
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

	// Preparing handler
	handler := http.HandlerFunc(mainHandle)

	// Filling structure with count parameter and awaited value
	requests := []CafeCountReq{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 100},
	}

	var cafeSlice []string

	var cafeCount int

	// Looping through all cities
	for city, _ := range cafeList {

		// Receiving maximum cafes count per a city
		maxCafeCount := len(cafeList[city])

		for _, v := range requests {

			cafeCount = 0

			response := httptest.NewRecorder()

			// Preparing URL
			url := "/cafe?city=" + city + "&count=" + strconv.Itoa(v.count)

			// Executing request
			req := httptest.NewRequest("GET", url, nil)

			handler.ServeHTTP(response, req)

			// Evaluating a received status code
			require.Equal(t, http.StatusOK, response.Code)

			// Getting cafes string from a request body
			cafeString := strings.TrimSpace(response.Body.String())

			if cafeString != "" {

				// Splitting a resulted string into a slice
				cafeSlice = strings.Split(cafeString, ",")

				// Getting a size of the slice
				cafeCount = len(cafeSlice)

			}

			// Performing a test validation for cafes count
			assert.Equal(t, min(v.wait, maxCafeCount), cafeCount)

		}
	}
}

func TestCafeSearch(t *testing.T) {

	// Preparing handler
	handler := http.HandlerFunc(mainHandle)

	// Filling structure with search parameter values

	requests := []CafeSearchReq{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	var cafeSlice []string

	var cafeCount int

	// Setting a city name for tests
	const city = "moscow"

	// Receiving maximum cafes count per a city
	maxCafeCount := len(cafeList[city])

	for _, v := range requests {

		cafeCount = 0

		response := httptest.NewRecorder()

		// Preparing URL
		url := "/cafe?city=" + city + "&search=" + v.search

		// Executing request
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		// Evaluating a received status code
		require.Equal(t, http.StatusOK, response.Code)

		// Getting cafes string from a request body
		cafeString := strings.TrimSpace(response.Body.String())

		if cafeString != "" {

			// Setting search string to a lower case
			searchStringLowerCase := strings.ToLower(v.search)

			// Splitting a resulted string into a slice
			cafeSlice = strings.Split(cafeString, ",")

			// Validating cafe names for search string
			for _, cafeName := range cafeSlice {

				// Setting cafe name to a lower case
				cafeNameLowerCase := strings.ToLower(cafeName)

				// Performing a validation, that a cafe name contains a search string (both in lower case)
				assert.True(t, strings.Contains(cafeNameLowerCase, searchStringLowerCase), true)

			}

			// Getting a size of the slice
			cafeCount = len(cafeSlice)

		}

		// Performing a test validation for cafes count
		assert.Equal(t, min(v.wantCount, maxCafeCount), cafeCount)

	}
}
