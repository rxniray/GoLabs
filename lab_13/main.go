package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/sum", sumHandler)

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}

// GET http://localhost:8080/sum?a=3&b=4
// {"result":7}
func sumHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		http.Error(w, "invalid parameter a", http.StatusBadRequest)
		return
	}
	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		http.Error(w, "invalid parameter b", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"result": sum(a, b)})
}

func sum(a, b float64) float64 {
	return a + b
}