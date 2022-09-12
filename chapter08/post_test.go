package chapter08

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type User struct {
	First string
	Last  string
}

func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(ioutil.Discard, r)
			_ = r.Close()
		}(r.Body)

		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			t.Error(err)
			http.Error(w, "Decode Failed", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}

}

func TestPostUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))
	defer ts.Close()

	response, err := http.Get(ts.URL)

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, actual status %d", http.StatusMethodNotAllowed, response.StatusCode)
	}

	buf := new(bytes.Buffer)
	u := User{First: "jung", Last: "ho"}
	err = json.NewEncoder(buf).Encode(&u)

	if err != nil {
		t.Fatal(err)
	}

	response, err = http.Post(ts.URL, "application/json", buf)

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d, actual status %d", http.StatusAccepted, response.StatusCode)
	}

	_ = response.Body.Close()
}

func TestMultipartPost(t *testing.T) {
	requestBody := new(bytes.Buffer)
	w := multipart.NewWriter(requestBody)

	m := map[string]string{
		"date":        time.Now().Format(time.RFC3339),
		"description": "Form values with attached files",
	}

	for k, v := range m {
		err := w.WriteField(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	files := []string{
		"./files/hello.txt",
		"./files/goodbye.txt",
	}

	for i, file := range files {
		filePart, err := w.CreateFormFile(fmt.Sprintf("file%d", i+1),
			filepath.Base(file))

		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(file)

		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(filePart, f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://httpbin.org/post", requestBody)

	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", w.FormDataContentType())

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, actual status %d", http.StatusOK, response.StatusCode)
	}

	t.Logf("\n%s", b)
}
