package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenURLHandle(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			"correct request",
			"https://practicum.yandex.ru/",
			want{http.StatusCreated},
		},
		{
			"empty request body",
			"",
			want{http.StatusBadRequest},
		},
		{
			"provided very short url",
			"http://i.io",
			want{http.StatusCreated},
		},
		{
			"provided not a url",
			"Hello World!",
			want{http.StatusBadRequest},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(test.body)))
			w := httptest.NewRecorder()
			ShortenURLHandle(w, r)

			res := w.Result()
			defer res.Body.Close()
			//resBody, err := io.ReadAll(res.Body)

			assertStatusCode(t, res.StatusCode, test.want.statusCode)
		})
	}
}

func TestGetURLHandle(t *testing.T) {
	type want struct {
		statusCode int
		locHeader  string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			"correct request",
			"EwHXdJfB",
			want{http.StatusTemporaryRedirect, "https://practicum.yandex.ru/"},
		},
	}

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/"+test.request, nil)
		r.SetPathValue("id", test.request)
		w := httptest.NewRecorder()
		GetURLHandle(w, r)

		res := w.Result()

		assertStatusCode(t, res.StatusCode, test.want.statusCode)

		loc := res.Header.Get("Location")

		if loc == "" {
			t.Error("no Location header was provided by server")
		}

		if loc != test.want.locHeader {
			t.Errorf("incorrect Location header: got %q; want %q", loc, test.want.locHeader)
		}
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("unexpected status code: got %d; want %d", got, want)
	}
}
