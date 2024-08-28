package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/halviet/shortener/internal/config"
	"github.com/halviet/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(test.body)))
			w := httptest.NewRecorder()
			ShortenURLHandle(&storage.Store{}, config.Config{BaseAddr: "http://localhost:8080"})(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}

func TestGetURLHandle(t *testing.T) {
	type want struct {
		statusCode int
		origin     string
	}

	tests := []struct {
		name  string
		urlID string
		want  want
	}{
		{
			"correct request",
			"qVYlmrQn",
			want{http.StatusTemporaryRedirect, "https://practicum.yandex.ru/"},
		},
	}

	store := storage.New()
	for _, test := range tests {
		store.SaveURL(storage.ShortURL{
			Origin: test.want.origin,
			Short:  test.urlID,
		})
	}

	r := chi.NewRouter()
	r.Get("/{id}", GetURLHandle(store))

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testGetRequest(t, ts, http.MethodGet, test.urlID)
			defer resp.Body.Close()
			assert.Equal(t, test.want.statusCode, resp.StatusCode)
			assert.Equal(t, test.want.origin, resp.Header.Get("Location"))
		})
	}
}

func testGetRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+"/"+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestJSONShortenURLHandle(t *testing.T) {
	type (
		want struct {
			statusCode int
		}
		payload struct {
			URL string `json:"url"`
		}
	)

	tests := []struct {
		name    string
		payload interface{}
		want    want
	}{
		{
			"correct request",
			payload{URL: "https://usjogdlhweujihogvnr.com/saf3/w4sddsg/dsg"},
			want{http.StatusCreated},
		},
		{
			"empty request body",
			payload{},
			want{http.StatusBadRequest},
		},
		{
			"int for url",
			struct {
				URL int `json:"url"`
			}{10},
			want{http.StatusBadRequest},
		},
	}

	endpoint := "/api/shorten"

	store := storage.New()
	cfg, err := config.New()
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", JSONShortenURLHandle(store, cfg))
	})

	ts := httptest.NewServer(r)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(test.payload)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, ts.URL+endpoint, &body)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			var res Response
			err = json.NewDecoder(resp.Body).Decode(&res)

			require.NoError(t, err)
		})
	}
}
