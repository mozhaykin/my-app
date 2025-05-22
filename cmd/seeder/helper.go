package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/doug-martin/goqu/v9/exp"

	"github.com/jackc/pgx/v5"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	"github.com/brianvoe/gofakeit/v7"
)

func Status() string {
	statuses := []string{
		"pending",
		"active",
		"inactive",
		"banned",
	}
	gofakeit.ShuffleAnySlice(statuses)

	return statuses[0]
}

func Contacts() []byte {
	contact := struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
	}{
		Phone: gofakeit.Phone(),
		Email: gofakeit.Email(),
	}

	data, err := json.Marshal(&contact)
	if err != nil {
		log.Fatalf("Contacts: json.Marshal: %v\n", err)
	}

	return data
}

func GenIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, batchSize)

	for range batchSize {
		ids = append(ids, uuid.New())
	}

	return ids
}

func Insert(ctx context.Context, conn *pgx.Conn, batch []any, table string) error {
	sql, _, err := goqu.Insert(table).Rows(batch...).ToSQL()
	if err != nil {
		return fmt.Errorf("goqu.Insert.ToSQL: %w", err)
	}

	_, err = conn.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("conn.Exec: %w", err)
	}

	return nil
}

func Tags() exp.LiteralExpression {
	tags := []exp.LiteralExpression{
		goqu.L("ARRAY ['home', 'primary']"),
		goqu.L("ARRAY ['work', 'secondary']"),
		goqu.L("ARRAY ['school', 'tertiary']"),
	}
	gofakeit.ShuffleAnySlice(tags)

	return tags[0]
}
