package httpapi

import (
	"encoding/json"
	"net/http"

	"pseudogo/internal/codegen"
	"pseudogo/internal/parser"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	var req ConvertRequest
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
		writeJSON(w, http.StatusInternalServerError, ConvertResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, ConvertResponse{GoCode: goCode})
}

func writeJSON(w http.ResponseWriter, status int, resp ConvertResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}