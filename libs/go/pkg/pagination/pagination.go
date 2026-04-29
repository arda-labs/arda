package pagination

import "strconv"

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Params struct {
	Limit  int
	Offset int
}

func Normalize(pageSize int, pageToken string) Params {
	return Params{
		Limit:  NormalizePageSize(pageSize),
		Offset: OffsetFromToken(pageToken),
	}
}

func NormalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

func OffsetFromToken(pageToken string) int {
	if pageToken == "" {
		return 0
	}
	offset, err := strconv.Atoi(pageToken)
	if err != nil || offset < 0 {
		return 0
	}
	return offset
}

func NextOffsetToken(itemCount, pageSize, offset int) string {
	if itemCount <= pageSize {
		return ""
	}
	return strconv.Itoa(offset + pageSize)
}
