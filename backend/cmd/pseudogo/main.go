package main

import (
	"fmt"
	"os"

	"pseudogo/internal/codegen"
	"pseudogo/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "pemakaian: pseudogo <file.pseudo> [-o output.go]")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputPath := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "-o" && i+1 < len(os.Args) {
			outputPath = os.Args[i+1]
			i++
		}
	}

	src, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gagal membaca file: %v\n", err)
		os.Exit(1)
	}

	file, err := parser.Parse(string(src), inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Kesalahan sintaks pseudocode:\n  %v\n", err)
		os.Exit(1)
	}

	goSrc, err := codegen.Generate(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Kesalahan saat menghasilkan kode Go:\n  %v\n", err)
		os.Exit(1)
	}

	if outputPath == "" {
		fmt.Print(goSrc)
		return
	}
	if err := os.WriteFile(outputPath, []byte(goSrc), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "gagal menulis file output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Berhasil: %s -> %s\n", inputPath, outputPath)
}
