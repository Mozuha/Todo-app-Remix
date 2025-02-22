package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetUIDFromCtxAndCast(ctx *gin.Context) (pgtype.UUID, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return pgtype.UUID{}, errors.New("userID not found in context")
	}

	userIDUuid, err := StringToUUID(userID.(string))
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("failed to convert userID to expected type: %v", err)
	}

	return userIDUuid, nil
}
