package pagination

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name      string
		pageSize  int
		pageToken string
		want      Params
	}{
		{name: "default page size", pageSize: 0, want: Params{Limit: DefaultPageSize, Offset: 0}},
		{name: "max page size", pageSize: MaxPageSize + 1, want: Params{Limit: MaxPageSize, Offset: 0}},
		{name: "valid token", pageSize: 10, pageToken: "30", want: Params{Limit: 10, Offset: 30}},
		{name: "invalid token", pageSize: 10, pageToken: "abc", want: Params{Limit: 10, Offset: 0}},
		{name: "negative token", pageSize: 10, pageToken: "-10", want: Params{Limit: 10, Offset: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.pageSize, tt.pageToken)
			if got != tt.want {
				t.Fatalf("Normalize() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestNextOffsetToken(t *testing.T) {
	tests := []struct {
		name      string
		itemCount int
		pageSize  int
		offset    int
		want      string
	}{
		{name: "no next page", itemCount: 10, pageSize: 10, offset: 0, want: ""},
		{name: "has next page", itemCount: 11, pageSize: 10, offset: 0, want: "10"},
		{name: "next page with offset", itemCount: 21, pageSize: 20, offset: 40, want: "60"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextOffsetToken(tt.itemCount, tt.pageSize, tt.offset)
			if got != tt.want {
				t.Fatalf("NextOffsetToken() = %q, want %q", got, tt.want)
			}
		})
	}
}
