package api

import (
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var star = []byte("\u2b50")

func getRequestID(r *http.Request) (uuid.UUID, error) {
	id := r.PathValue("id")
	return uuid.Parse(id)
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
