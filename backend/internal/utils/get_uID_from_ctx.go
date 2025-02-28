package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetUIDFromCtxAndCast(ctx *gin.Context) (pgtype.UUID, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return pgtype.UUID{}, ErrUIDNotFoundInCtx
	}

	userIDUuid, err := StringToUUID(userID.(string))
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("failed to convert userID to expected type: %v", err)
	}

	return userIDUuid, nil
}

func GetUIDFromCtxAndCreateRespUponErr(ctx *gin.Context) (pgtype.UUID, error) {
	userIDUuid, err := GetUIDFromCtxAndCast(ctx)

	if err != nil {
		log.Println(err.Error())

		if err == ErrUIDNotFoundInCtx {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": MsgUIDNotFoundInCtx})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": MsgInternalServerErr})
		}
	}

	return userIDUuid, err
}
