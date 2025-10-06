package sorting

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetKeyColumn(t *testing.T) {
	tests := []struct {
		name string
		line string
		opts SortOptions
		want string
	}{
		{
			name: "extract 2nd column using tab separator",
			line: "first\tsecond\tthird",
			opts: SortOptions{Key: 2, Separator: '\t'},
			want: "second",
		},
		{
			name: "key out of range with tab separator",
			line: "onlyone",
			opts: SortOptions{Key: 2, Separator: '\t'},
			want: "onlyone",
		},
		{
			name: "key = 0 returns whole line",
			line: "hello wombat",
			opts: SortOptions{Key: 0, Separator: ' '},
			want: "hello wombat",
		},
		{
			name: "extract 3rd column with space separator",
			line: "a b c d",
			opts: SortOptions{Key: 3, Separator: ' '},
			want: "c",
		},
		{
			name: "multiple spaces treated as one when separator is 0",
			line: "a     b    c",
			opts: SortOptions{Key: 2, Separator: 0},
			want: "b",
		},
		{
			name: "separator is 0 (default), fields split by space/tab",
			line: "alpha\tbeta gamma",
			opts: SortOptions{Key: 2, Separator: 0},
			want: "beta",
		},
		{
			name: "separator explicitly set to space",
			line: "x y z",
			opts: SortOptions{Key: 2, Separator: ' '},
			want: "y",
		},
		{
			name: "key exceeds number of fields with space separator",
			line: "one two",
			opts: SortOptions{Key: 4, Separator: ' '},
			want: "one two",
		},
		{
			name: "empty line should return empty string",
			line: "",
			opts: SortOptions{Key: 1, Separator: '\t'},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getKeyColumn(tt.line, tt.opts)
			if got != tt.want {
				t.Errorf("getKeyColumn() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestCompareKeys(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		opts SortOptions
		want bool
	}{

		{"numeric ascending", "10", "2", SortOptions{NumericSort: true}, false},
		{"numeric descending", "10", "2", SortOptions{NumericSort: true, Reverse: true}, false},

		{"ignore blanks", "abc  ", "abc", SortOptions{IgnoreBlanks: true}, false},

		{"month order jan < feb", "Jan", "Feb", SortOptions{Month: true}, true},
		{"month order dec > nov", "Dec", "Nov", SortOptions{Month: true}, false},

		{"human readable K < M", "1K", "1M", SortOptions{Human: true}, true},
		{"human readable G > M", "2G", "1M", SortOptions{Human: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyA := getKeyColumn(tt.a, tt.opts)
			keyB := getKeyColumn(tt.b, tt.opts)
			got := compareKeys(tt.a, tt.b, keyA, keyB, tt.opts)
			if got != tt.want {
				t.Errorf("compareKeys(%q, %q) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}

}

func TestParseHuman(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		ok    bool
	}{
		{"1K", 1 << 10, true},
		{"2M", 2 << 20, true},
		{"3.5G", 3.5 * (1 << 30), true},
		{"3.5Gв", 3.5 * (1 << 30), true},
		{"abc", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ok, got := parseHuman(tt.input)
			if ok != tt.ok || got != tt.want {
				t.Errorf("tryParseHuman(%q) = (%v, %v); want (%v, %v)", tt.input, ok, got, tt.ok, tt.want)
			}
		})
	}
}

func TestCheckSorted(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		opts    SortOptions
		wantErr bool
	}{
		{
			name:    "empty input",
			lines:   []string{},
			opts:    SortOptions{},
			wantErr: false,
		},
		{
			name:    "lexicographically sorted lines",
			lines:   []string{"a", "b", "c"},
			opts:    SortOptions{},
			wantErr: false,
		},
		{
			name:    "lexicographically unsorted lines",
			lines:   []string{"a", "c", "b"},
			opts:    SortOptions{},
			wantErr: true,
		},
		{
			name:    "numeric sort (-n) correct",
			lines:   []string{"1", "2", "10"},
			opts:    SortOptions{NumericSort: true},
			wantErr: false,
		},
		{
			name:    "numeric sort (-n) incorrect",
			lines:   []string{"1", "10", "2"},
			opts:    SortOptions{NumericSort: true},
			wantErr: true,
		},
		{
			name:    "reverse sort (-r) correct",
			lines:   []string{"c", "b", "a"},
			opts:    SortOptions{Reverse: true},
			wantErr: false,
		},
		{
			name:    "reverse sort (-r) incorrect",
			lines:   []string{"a", "b", "c"},
			opts:    SortOptions{Reverse: true},
			wantErr: true,
		},
		{
			name:    "human-readable sort (-h) correct",
			lines:   []string{"1K", "1M", "1G"},
			opts:    SortOptions{Human: true},
			wantErr: false,
		},
		{
			name:    "human-readable sort (-h) incorrect",
			lines:   []string{"1K", "1G", "1M"},
			opts:    SortOptions{Human: true},
			wantErr: true,
		},
		{
			name:    "month sort (-M) correct",
			lines:   []string{"Jan", "Feb", "Mar"},
			opts:    SortOptions{Month: true},
			wantErr: false,
		},
		{
			name:    "month sort (-M) incorrect",
			lines:   []string{"Jan", "Mar", "Feb"},
			opts:    SortOptions{Month: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSorted(tt.lines, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkSorted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func runSortTest(t *testing.T, opts SortOptions, unixArgs []string) {
	t.Helper()

	testFile := "test/test.txt"
	testDir := "test"

	tmpOut := filepath.Join(testDir, "tmp_out.txt")
	tmpUnix := filepath.Join(testDir, "tmp_unix.txt")

	outFile, err := os.Create(tmpOut)
	if err != nil {
		t.Fatalf("failed to create tmp output file: %v", err)
	}
	oldStdout := os.Stdout
	os.Stdout = outFile
	if err := SortLines(opts); err != nil {
		t.Fatalf("SortLines failed: %v", err)
	}
	outFile.Close()
	os.Stdout = oldStdout

	cmd := exec.Command("sort", append(unixArgs, testFile)...)
	unixOut, err := os.Create(tmpUnix)
	if err != nil {
		t.Fatalf("failed to create unix output file: %v", err)
	}
	cmd.Stdout = unixOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("unix sort failed: %v", err)
	}
	unixOut.Close()

	gotLines, err := os.ReadFile(tmpOut)
	if err != nil {
		t.Fatalf("failed to read tmp output file: %v", err)
	}
	wantLines, err := os.ReadFile(tmpUnix)
	if err != nil {
		t.Fatalf("failed to read unix output file: %v", err)
	}

	got := strings.Split(strings.TrimSpace(string(gotLines)), "\n")
	want := strings.Split(strings.TrimSpace(string(wantLines)), "\n")

	if len(got) != len(want) {
		t.Fatalf("line count mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}

	// cleanup
	if !t.Failed() {
		_ = os.Remove(tmpOut)
		_ = os.Remove(tmpUnix)
	} else {
		t.Logf("test failed, output files left in %s", testDir)
	}
}

func TestSortLines_CompareWithUnixSort_Files(t *testing.T) {
	tests := []struct {
		name     string
		opts     SortOptions
		unixArgs []string
	}{
		{
			name:     "simple",
			opts:     SortOptions{Files: []string{"test/test.txt"}},
			unixArgs: []string{},
		}, {
			name:     "reverse",
			opts:     SortOptions{Files: []string{"test/test.txt"}, Reverse: true},
			unixArgs: []string{"-r"},
		},
		{
			name:     "numeric",
			opts:     SortOptions{Files: []string{"test/test.txt"}, NumericSort: true},
			unixArgs: []string{"-n"},
		}, {
			name:     "unique",
			opts:     SortOptions{Files: []string{"test/test.txt"}, Unique: true},
			unixArgs: []string{"-u"},
		}, {
			name:     "month",
			opts:     SortOptions{Files: []string{"test/test.txt"}, Month: true},
			unixArgs: []string{"-M"},
		}, {
			name:     "humanReadble",
			opts:     SortOptions{Files: []string{"test/test.txt"}, Human: true},
			unixArgs: []string{"-h"},
		}, {
			name:     "ignoreBlanks",
			opts:     SortOptions{Files: []string{"test/test.txt"}, IgnoreBlanks: true},
			unixArgs: []string{"-b"},
		}, {
			name:     "key",
			opts:     SortOptions{Files: []string{"test/test.txt"}, Key: 2},
			unixArgs: []string{"-k2"},
		}, {
			name: "reverse_unique",
			opts: SortOptions{
				Files:   []string{"test/test.txt"},
				Reverse: true,
				Unique:  true,
			},
			unixArgs: []string{"-r", "-u"},
		}, {
			name: "numeric_reverse",
			opts: SortOptions{
				Files:       []string{"test/test.txt"},
				NumericSort: true,
				Reverse:     true,
			},
			unixArgs: []string{"-n", "-r"},
		},
		{
			name: "numeric_unique",
			opts: SortOptions{
				Files:       []string{"test/test.txt"},
				NumericSort: true,
				Unique:      true,
			},
			unixArgs: []string{"-n", "-u"},
		},
		{
			name: "month_ignoreBlanks",
			opts: SortOptions{
				Files:        []string{"test/test.txt"},
				Month:        true,
				IgnoreBlanks: true,
			},
			unixArgs: []string{"-M", "-b"},
		},
		{
			name: "human_reverse",
			opts: SortOptions{
				Files:   []string{"test/test.txt"},
				Human:   true,
				Reverse: true,
			},
			unixArgs: []string{"-h", "-r"},
		},
		{
			name: "human_unique_ignoreBlanks",
			opts: SortOptions{
				Files:        []string{"test/test.txt"},
				Human:        true,
				Unique:       true,
				IgnoreBlanks: true,
			},
			unixArgs: []string{"-h", "-u", "-b"},
		},
		{
			name: "key_numeric",
			opts: SortOptions{
				Files:       []string{"test/test.txt"},
				Key:         2,
				NumericSort: true,
			},
			unixArgs: []string{"-k2", "-n"},
		},
		{
			name: "key_human_reverse",
			opts: SortOptions{
				Files:   []string{"test/test.txt"},
				Key:     2,
				Human:   true,
				Reverse: true,
			},
			unixArgs: []string{"-k2", "-h", "-r"},
		},
		{
			name: "key_month_ignoreBlanks",
			opts: SortOptions{
				Files:        []string{"test/test.txt"},
				Key:          2,
				Month:        true,
				IgnoreBlanks: true,
			},
			unixArgs: []string{"-k2", "-M", "-b"},
		},
		{
			name: "everything_combined",
			opts: SortOptions{
				Files:        []string{"test/test.txt"},
				Key:          2,
				Human:        true,
				Reverse:      true,
				Unique:       true,
				IgnoreBlanks: true,
			},
			unixArgs: []string{"-k2", "-h", "-r", "-u", "-b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runSortTest(t, tc.opts, tc.unixArgs)
		})
	}
}
