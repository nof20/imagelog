package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// 100% transparent 1 px PNG file
// See read.go
const transpng = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR4nGP6zwAAAgcBApocMXEAAAAASUVORK5CYII="

func main() {
	http.HandleFunc("/imagelog.png", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Save request data to sheet
	ctx := appengine.NewContext(r)
	saveToSheet(ctx, r)

	// Return 1 px PNG
	dat, _ := base64.StdEncoding.DecodeString(transpng)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	fmt.Fprint(w, dat)
}

func saveToSheet(ctx context.Context, r *http.Request) {
	// Get data from request headers.
	headers := []string{
		"User-Agent",
		"X-Appengine-Country",
		"X-AppEngine-Region",
		"X-AppEngine-City",
		"X-AppEngine-CityLatLong",
		"X-Appengine-Remote-Addr",
		"X-Appengine-User-Organization",
		"X-Appengine-User-Email",
		"Referer",
	}
	length := 2 + len(headers)
	results := make([]interface{}, length)
	for i, h := range headers {
		results[i] = r.Header.Get(h)
	}

	// Add date/time.
	results[length-2] = time.Now().UTC().Format(time.RFC3339)
	results[length-1] = r.URL.Path

	// Save data to Google Sheet.
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, sheets.SpreadsheetsScope),
			Base: &urlfetch.Transport{
				Context: ctx,
			},
		},
	}
	svc, err := sheets.New(client)
	if err != nil {
		log.Errorf(ctx, "Unable to retrieve Sheets client: %v", err)
	}
	spreadsheetID := os.Getenv("SHEET_ID")
	sheet := "Log"
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, results)
	_, err = svc.Spreadsheets.Values.Append(spreadsheetID, sheet, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Errorf(ctx, "Unable to write data to sheet. %v", err)
	}
}
