package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}

	src := uuid.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func StringToUUID(s string) (pgtype.UUID, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Scan(s); err != nil {
		return pgtype.UUID{}, err
	}

	return uuid, nil
}
