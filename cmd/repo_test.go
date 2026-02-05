package cmd

import (
	"slices"
	"testing"
)

func TestParseOverwrittenFilePaths_CheckoutUntracked(t *testing.T) {
	out := `error: The following untracked working tree files would be overwritten by checkout:
	.bashrc
	.config/nvim/init.vim
Please move or remove them before you switch branches.
Aborting`

	got := parseOverwrittenFilePaths(out)
	want := []string{".bashrc", ".config/nvim/init.vim"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestParseOverwrittenFilePaths_MergeLocalChanges(t *testing.T) {
	out := `error: Your local changes to the following files would be overwritten by merge:
	.zshrc
Please commit your changes or stash them before you merge.
Aborting`

	got := parseOverwrittenFilePaths(out)
	want := []string{".zshrc"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestParseOverwrittenFilePaths_NoConflicts(t *testing.T) {
	out := "Already up to date.\n"
	got := parseOverwrittenFilePaths(out)
	if len(got) != 0 {
		t.Fatalf("got %v, want empty", got)
	}
}

func TestSafeRelPath(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "normal", in: ".zshrc", want: ".zshrc"},
		{name: "traversal", in: "../evil", wantErr: true},
		{name: "absolute", in: "/etc/passwd", wantErr: true},
		{name: "empty", in: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := safeRelPath(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
