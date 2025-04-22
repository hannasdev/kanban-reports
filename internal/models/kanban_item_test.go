// internal/models/kanban_item_test.go
package models

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "Valid timestamp",
			timeStr:  "2024/05/07 03:49:34",
			wantTime: time.Date(2024, 5, 7, 3, 49, 34, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Empty string",
			timeStr:  "",
			wantTime: time.Time{},
			wantErr:  false,
		},
		{
			name:     "Invalid format",
			timeStr:  "07-05-2024",
			wantTime: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := ParseTime(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !gotTime.Equal(tt.wantTime) {
				t.Errorf("ParseTime() = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}

func TestParseOwners(t *testing.T) {
	tests := []struct {
		name      string
		ownersStr string
		want      []string
	}{
		{
			name:      "Comma separated",
			ownersStr: "john.doe@example.com,jane.smith@example.com",
			want:      []string{"john.doe@example.com", "jane.smith@example.com"},
		},
		{
			name:      "Single owner",
			ownersStr: "john.doe@example.com",
			want:      []string{"john.doe@example.com"},
		},
		{
			name:      "Empty string",
			ownersStr: "",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOwners(tt.ownersStr)
			if len(got) != len(tt.want) {
				t.Errorf("ParseOwners() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ParseOwners()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}