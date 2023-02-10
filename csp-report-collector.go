package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// CSPReport represents the structure of a CSP report
type CSPReport struct {
	DocumentURI        string `json:"document-uri"`
	Referrer           string `json:"referrer"`
	BlockedURI         string `json:"blocked-uri"`
	ViolatedDirective  string `json:"violated-directive"`
	EffectiveDirective string `json:"effective-directive"`
	OriginalPolicy     string `json:"original-policy"`
	Disposition        string `json:"disposition"`
	StatusCode         int    `json:"status-code"`
}

type Config struct {
	Datasrouce string `json:"datasource"`
}

func getConfig() (*Config, error) {
	file, fileErr := os.ReadFile(os.Args[1])
	if fileErr != nil {
		return nil, fileErr
	}
	var config Config
	err := json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, configErr := getConfig()
	if configErr != nil {
		fmt.Println("Could not parse config")
		return
	}
	db, err := sql.Open("mysql", config.Datasrouce)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/csp-report", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "POST" {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Report CSPReport `json:"csp-report"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := db.Exec("INSERT INTO csp_reports (document_uri, referrer, blocked_uri, violated_directive, effective_directive, original_policy, disposition, status_code) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", payload.Report.DocumentURI, payload.Report.Referrer, payload.Report.BlockedURI, payload.Report.ViolatedDirective, payload.Report.EffectiveDirective, payload.Report.OriginalPolicy, payload.Report.Disposition, payload.Report.StatusCode); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	if err := http.ListenAndServe(":3010", nil); err != nil {
		panic(err)
	}
}
