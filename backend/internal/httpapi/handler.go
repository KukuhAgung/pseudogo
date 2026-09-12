package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"pseudogo/internal/codegen"
	"pseudogo/internal/parser"
	"pseudogo/internal/sandbox"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
	r.Body = http.MaxBytesReader(w, r.Body, 100*1024)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: "request body tidak valid"})
		return
	}

	file, err := parser.Parse(req.Pseudocode, "input")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ConvertResponse{Error: err.Error()})
		return
	}

	goCode, err := codegen.Generate(file)
	if err != nil {
		log.Printf("codegen error: %v", err)
		writeJSON(w, http.StatusInternalServerError, ConvertResponse{Error: "Terjadi kesalahan internal saat memproses kode. Coba lagi beberapa saat lagi."})
		return
	}

	writeJSON(w, http.StatusOK, ConvertResponse{GoCode: goCode})
}

func writeJSON(w http.ResponseWriter, status int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 100*1024)

	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, RunResponse{Error: "request body tidak valid"})
		return
	}

	file, err := parser.Parse(req.Pseudocode, "input")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, RunResponse{Error: err.Error()})
		return
	}

	goCode, err := codegen.Generate(file)
	if err != nil {
		log.Printf("codegen error: %v", err)
		writeJSON(w, http.StatusInternalServerError, RunResponse{Error: "Terjadi kesalahan internal saat memproses kode. Coba lagi beberapa saat lagi."})
		return
	}

	result, err := sandbox.RunGoCode(goCode, req.Input)
	if err != nil {
		log.Printf("sandbox error: %v", err)
		writeJSON(w, http.StatusInternalServerError, RunResponse{Error: "Terjadi kesalahan pada server saat menjalankan program. Coba lagi beberapa saat lagi."})
		return
	}

	writeJSON(w, http.StatusOK, RunResponse{Output: result.Stdout, Stderr: result.Stderr, TimedOut: result.TimedOut})
}
