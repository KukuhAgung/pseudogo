package sandbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const pistonURL = "http://localhost:2000/api/v2/execute"

type pistonFile struct {
	Content string `json:"content"`
}

type pistonRequest struct {
	Language       string       `json:"language"`
	Version        string       `json:"version"`
	Files          []pistonFile `json:"files"`
	Stdin          string       `json:"stdin"`
	RunTimeout     int          `json:"run_timeout"`
	RunCPUTime     int          `json:"run_cpu_time"`
	RunMemoryLimit int          `json:"run_memory_limit"`
}

type pistonStage struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Code     int    `json:"code"`
	Status   string `json:"status"`
	CPUTime  int    `json:"cpu_time"`
	WallTime int    `json:"wall_time"`
}

type pistonResponse struct {
	Run     pistonStage `json:"run"`
	Compile pistonStage `json:"compile"`
	Message string      `json:"message"`
}

type RunResult struct {
	Stdout   string
	Stderr   string
	TimedOut bool
	WallTime int
}

func RunGoCode(goSource string, stdinInput string) (*RunResult, error) {
	body, _ := json.Marshal(pistonRequest{
		Language:       "go",
		Version:        "*",
		Files:          []pistonFile{{Content: goSource}},
		Stdin:          stdinInput,
		RunTimeout:     25000,
		RunCPUTime:     15000,
		RunMemoryLimit: 67108864,
	})
	
	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Post(pistonURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi Piston: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var result pistonResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("gagal baca respons Piston: %w", err)
	}
	if result.Message != "" {
		return nil, fmt.Errorf("Piston menolak request: %s", result.Message)
	}
	if result.Compile.Code != 0 {
		return &RunResult{Stderr: "Gagal compile:\n" + result.Compile.Stderr}, nil
	}

	return &RunResult{
		Stdout:   result.Run.Stdout,
		Stderr:   result.Run.Stderr,
		TimedOut: result.Run.Status == "TO",
	}, nil
}
