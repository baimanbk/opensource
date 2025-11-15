package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/avg", avgHandler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type result struct {
	Sum     int     `json:"sum,omitempty"`
	Average float64 `json:"average,omitempty"`
	Error   string  `json:"error,omitempty"`
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("nums")
	if q == "" {
		writeJSON(w, result{Error: "missing nums query param, e.g. ?nums=1,2,3"})
		return
	}
	nums, err := parseNums(q)
	if err != nil {
		writeJSON(w, result{Error: err.Error()})
		return
	}
	s := 0
	for _, n := range nums {
		s += n
	}
	writeJSON(w, result{Sum: s})
}

func avgHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("nums")
	if q == "" {
		writeJSON(w, result{Error: "missing nums query param, e.g. ?nums=1,2,3"})
		return
	}
	nums, err := parseNums(q)
	if err != nil {
		writeJSON(w, result{Error: err.Error()})
		return
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	avg := float64(sum) / float64(len(nums))
	writeJSON(w, result{Average: avg})
}

func parseNums(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", p, err)
		}
		// append parsed number to the slice instead of assigning by index
		nums = append(nums, n)
	}
	return nums, nil
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
