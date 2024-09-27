package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const templatePath = "templates"

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
