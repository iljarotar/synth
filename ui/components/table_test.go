package components

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_truncateRows(t *testing.T) {
	tests := []struct {
		name         string
		rows         []Row
		selected     int
		height       int
		wantRow      []Row
		wantSelected int
	}{
		{
			name: "height big enough no truncation necessary",
			rows: []Row{
				{
					Columns: []string{"1"},
				},
			},
			selected: 0,
			height:   10,
			wantRow: []Row{
				{
					Columns: []string{"1"},
				},
			},
		},
		{
			name: "middle selected",
			rows: []Row{
				{
					Columns: []string{"1"},
				},
				{
					Columns: []string{"2"},
				},
				{
					Columns: []string{"3"},
				},
				{
					Columns: []string{"4"},
				},
				{
					Columns: []string{"5"},
				},
			},
			selected: 2,
			height:   3,
			wantRow: []Row{
				{
					Columns: []string{"2"},
				},
				{
					Columns: []string{"3"},
				},
				{
					Columns: []string{"4"},
				},
			},
			wantSelected: 1,
		},
		{
			name: "lower middle selected even number",
			rows: []Row{
				{
					Columns: []string{"1"},
				},
				{
					Columns: []string{"2"},
				},
				{
					Columns: []string{"3"},
				},
				{
					Columns: []string{"4"},
				},
				{
					Columns: []string{"5"},
				},
				{
					Columns: []string{"6"},
				},
			},
			selected: 2,
			height:   3,
			wantRow: []Row{
				{
					Columns: []string{"2"},
				},
				{
					Columns: []string{"3"},
				},
				{
					Columns: []string{"4"},
				},
			},
			wantSelected: 1,
		},
		{
			name: "0 selected",
			rows: []Row{
				{
					Columns: []string{"1"},
				},
				{
					Columns: []string{"2"},
				},
				{
					Columns: []string{"3"},
				},
				{
					Columns: []string{"4"},
				},
				{
					Columns: []string{"5"},
				},
				{
					Columns: []string{"6"},
				},
			},
			selected: 0,
			height:   4,
			wantRow: []Row{
				{
					Columns: []string{"1"},
				},
				{
					Columns: []string{"2"},
				},
				{
					Columns: []string{"3"},
				},
				{
					Columns: []string{"4"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, gotSelected := truncateRows(tt.rows, tt.selected, tt.height)
			if diff := cmp.Diff(tt.wantRow, gotRows); diff != "" {
				t.Errorf("truncateRows() diff = %v", diff)
			}
			if gotSelected != tt.wantSelected {
				t.Errorf("truncateRows() got selected = %v, want = %v", gotSelected, tt.wantSelected)
			}
		})
	}
}
