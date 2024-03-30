package helper

import (
	"fmt"
	"strconv"
)

const (
	decimal = 10
	bitSize = 64
)

func ParseUint64(idParam string) (uint64, error) {
	id, err := strconv.ParseUint(idParam, decimal, bitSize)
	if err != nil {
		return 0, fmt.Errorf("failed to parse id: %w", err)
	}

	return id, nil
}
