package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var star = []byte("\u2b50")

func getTemplate(fname string, funcs template.FuncMap) *template.Template {
	fp := filepath.Join(templatePath, fname)
	return template.Must(
		template.New(fname).Funcs(funcs).ParseFiles(fp),
	)
}

func getRequestID(r *http.Request) (uuid.UUID, error) {
	id := r.PathValue("id")
	return uuid.Parse(id)
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

func uuidToPgType(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func timeToPgType(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

func intToPgType(i int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(i), Valid: true}
}

func getDifficultyString(difficulty int) string {
	r, _ := utf8.DecodeRune(star)
	return strings.Repeat(string(r), difficulty)
}

func difficultyStringToInt(s string) int {
	return strings.Count(s, string(star))
}

func fetchRecord[T any](client *http.Client, endpoint string, queryParams queryParams) (T, error) {
	var record T

	url := baseURL + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return record, err
	}

	if queryParams != nil {
		req.URL.RawQuery = buildQuery(queryParams)
	}

	resp, err := client.Do(req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		msg := fmt.Sprintf("%d %s %s", resp.StatusCode, resp.Request.Method, endpoint)
		return record, errors.New(msg)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return record, err
	}

	if err := json.Unmarshal(b, &record); err != nil {
		return record, err
	}

	return record, nil
}

type queryParams map[string]string

func buildQuery(params queryParams) string {
	q := url.Values{}

	for k, v := range params {
		q.Add(k, v)
	}

	return q.Encode()
}
