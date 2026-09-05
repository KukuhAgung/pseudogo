package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"pseudogo/internal/codegen"
	"pseudogo/internal/parser"
)

// convert parses+generates a .pseudo file and returns the formatted Go source.
func convert(t *testing.T, path string) string {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("gagal baca %s: %v", path, err)
	}
	file, err := parser.Parse(string(src), path)
	if err != nil {
		t.Fatalf("parse error di %s: %v", path, err)
	}
	goSrc, err := codegen.Generate(file)
	if err != nil {
		t.Fatalf("codegen error di %s: %v", path, err)
	}
	return goSrc
}

// buildAndRun writes goSrc to a temp file, compiles it, runs it with stdin,
// and returns combined stdout.
func buildAndRun(t *testing.T, goSrc string, stdin string) string {
	t.Helper()
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(srcPath, []byte(goSrc), 0644); err != nil {
		t.Fatalf("gagal tulis source: %v", err)
	}
	binPath := filepath.Join(dir, "prog")
	build := exec.Command("go", "build", "-o", binPath, srcPath)
	var buildErr bytes.Buffer
	build.Stderr = &buildErr
	if err := build.Run(); err != nil {
		t.Fatalf("go build gagal: %v\n%s\n--- generated source ---\n%s", err, buildErr.String(), goSrc)
	}
	run := exec.Command(binPath)
	run.Stdin = strings.NewReader(stdin)
	var out, runErr bytes.Buffer
	run.Stdout = &out
	run.Stderr = &runErr
	if err := run.Run(); err != nil {
		t.Fatalf("program gagal jalan: %v\nstderr: %s", err, runErr.String())
	}
	return out.String()
}

// TestAllTestdataCompiles is a smoke test: every .pseudo file under
// testdata/ must parse, generate valid Go, and compile successfully.
func TestAllTestdataCompiles(t *testing.T) {
	files, err := filepath.Glob("../testdata/*.pseudo")
	if err != nil || len(files) == 0 {
		t.Fatalf("tidak ada file testdata ditemukan: %v", err)
	}
	for _, f := range files {
		f := f
		t.Run(filepath.Base(f), func(t *testing.T) {
			goSrc := convert(t, f)
			dir := t.TempDir()
			srcPath := filepath.Join(dir, "main.go")
			os.WriteFile(srcPath, []byte(goSrc), 0644)
			binPath := filepath.Join(dir, "prog")
			build := exec.Command("go", "build", "-o", binPath, srcPath)
			var buildErr bytes.Buffer
			build.Stderr = &buildErr
			if err := build.Run(); err != nil {
				t.Fatalf("go build gagal: %v\n%s\n--- generated source ---\n%s", err, buildErr.String(), goSrc)
			}
		})
	}
}

func TestSequentialSearch(t *testing.T) {
	goSrc := convert(t, "../testdata/sequential_search.pseudo")
	out := buildAndRun(t, goSrc, "")
	want := "Indeks ditemukan (1-based): 4\nHasil pencarian tidak ada: -1\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}

func TestBinarySearch(t *testing.T) {
	goSrc := convert(t, "../testdata/binary_search.pseudo")
	out := buildAndRun(t, goSrc, "")
	want := "Indeks 23 (1-based): 6\nHasil pencarian tidak ada: -1\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}

func TestSelectionSort(t *testing.T) {
	goSrc := convert(t, "../testdata/selection_sort.pseudo")
	out := buildAndRun(t, goSrc, "")
	want := "data[1] = 11\ndata[2] = 12\ndata[3] = 22\ndata[4] = 25\ndata[5] = 64\ndata[6] = 90\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}

func TestInsertionSort(t *testing.T) {
	goSrc := convert(t, "../testdata/insertion_sort.pseudo")
	out := buildAndRun(t, goSrc, "")
	want := "data[1] = 1\ndata[2] = 3\ndata[3] = 4\ndata[4] = 5\ndata[5] = 7\ndata[6] = 9\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}

func TestMiscFeatures(t *testing.T) {
	goSrc := convert(t, "../testdata/misc_features.pseudo")
	out := buildAndRun(t, goSrc, "")
	want := "Nama mahasiswa: Budi\nIndeks A\nIndeks huruf: A\nNilai PHI: 3.14\n6 genap? true\n7 genap? false\nHitungan ke: 1\nHitungan ke: 2\nHitungan ke: 3\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}

func TestHitungLuasAndInput(t *testing.T) {
	goSrc := convert(t, "../testdata/hitung_luas.pseudo")
	out := buildAndRun(t, goSrc, "5 8 Budi 20 W\n")
	want := "Luas = 40\nHalo Budi, umur 20, inisial W\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}

func TestNestedOutParam(t *testing.T) {
	goSrc := convert(t, "../testdata/nested_out.pseudo")
	out := buildAndRun(t, goSrc, "")
	want := "Hasil: 70\n"
	if out != want {
		t.Errorf("output tidak sesuai.\ngot:  %q\nwant: %q", out, want)
	}
}
