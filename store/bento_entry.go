package store

import (
	"fmt"
	"strings"
	"time"
)

// BentoEntry represents an entry (secret) of a already preprared bento in the database.
type BentoEntry struct {
	Id        int64
	Name      string
	Value     string
	BentoId   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBentoEntry just creates a new BentoEntry struct and it DOES NOT saves it in the database.
// This allows any procedure to create multiple BentoEntry if needed and save them in a batch.
func NewBentoEntry(name, value, bentoId string) BentoEntry {
	return BentoEntry{Name: name, Value: value, BentoId: bentoId}
}

// SaveBentoEntryBatch will save a batch of BentoEntry at once. The entries will be updated in place
// with the corresponding id, created_at and updated_at from the database.
func SaveBentoEntryBatch(entries []BentoEntry) error {
	if len(entries) == 0 {
		return nil
	}
	// need to build a string with all the values first
	var builder strings.Builder
	var values []any
	c := 1
	for i, e := range entries {
		builder.WriteString(fmt.Sprintf("($%d, $%d, $%d)", c, c+1, c+2))
		c += 3
		if i < len(entries)-1 {
			builder.WriteString(",")
		}
		values = append(values, e.Name, e.Value, e.BentoId)
	}
	sqlStr := fmt.Sprintf("INSERT INTO bento_entries (name, value, bento_id) VALUES %s RETURNING id, created_at, updated_at;", builder.String())
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&entries[i].Id, &entries[i].CreatedAt, &entries[i].UpdatedAt)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

// GetEntriesForBento will get all the entries for a given bento.
func GetEntriesForBento(bentoId string) ([]BentoEntry, error) {
	rows, err := db.Query("SELECT id, name, value, bento_id, created_at, updated_at FROM bento_entries WHERE bento_id = $1;", bentoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := []BentoEntry{}
	for rows.Next() {
		entry := BentoEntry{}
		err = rows.Scan(
			&entry.Id,
			&entry.Name,
			&entry.Value,
			&entry.BentoId,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
