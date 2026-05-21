package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// defaultOffset это дефолтное значение для смещения при пагинации.
	defaultOffset = 0

	// defaultLimit это дефолтное значения для лимита при пагинации.
	defaultLimit = 20

	// maxLimit это максимальное значение для лимита при пагинации.
	maxLimit = 100
)

// getUintParam извлекает и валидирует целочисленный ненулевой параметр из URL, например id.
func getUUIDParam(c *gin.Context, name string) (uuid.UUID, error) {
	var UUIDParam uuid.UUID

	stringUUID := c.Param(name)
	if stringUUID == "" {
		return UUIDParam, nil
	}

	UUIDParam, err := uuid.Parse(stringUUID)
	return UUIDParam, err
}

// getPaginationParams извлекает параметры пагинации из HTTP-запроса
func getPaginationParams(c *gin.Context) (uint64, uint64, error) {
	offset := defaultOffset
	limit := defaultLimit

	if rawOffset := c.Query("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil || parsedOffset < 0 {
			return 0, 0, errors.New("invalid offset")
		}

		offset = parsedOffset
	}

	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 {
			return 0, 0, errors.New("invalid limit")
		}

		if parsedLimit > maxLimit {
			parsedLimit = maxLimit
		}

		limit = parsedLimit
	}

	return uint64(limit), uint64(offset), nil
}
