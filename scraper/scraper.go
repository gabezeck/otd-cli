package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

type Event struct {
	Year string
	Text string
}

type Birthday struct {
	Name string
	Year string
}

type OTDData struct {
	Date      string
	Events    []Event
	Birthdays []Birthday
}

type cacheFile struct {
	Date string  // "YYYY-MM-DD" of when this was fetched
	Data OTDData
}

func cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "otd-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func loadCache() (*OTDData, bool) {
	dir, err := cacheDir()
	if err != nil {
		return nil, false
	}
	path := filepath.Join(dir, "today.json")

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var cf cacheFile
	if err := json.Unmarshal(b, &cf); err != nil {
		return nil, false
	}

	today := time.Now().Format("2006-01-02")
	if cf.Date != today {
		return nil, false
	}

	return &cf.Data, true
}

func writeCache(data *OTDData) {
	dir, err := cacheDir()
	if err != nil {
		return
	}
	cf := cacheFile{
		Date: time.Now().Format("2006-01-02"),
		Data: *data,
	}
	b, err := json.Marshal(cf)
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(dir, "today.json"), b, 0644)
}

func Scrape() (*OTDData, error) {
	if data, ok := loadCache(); ok {
		return data, nil
	}

	data, err := fetch()
	if err != nil {
		return nil, err
	}

	writeCache(data)
	return data, nil
}

func fetch() (*OTDData, error) {
	// Load .env if present; ignore error (file may not exist)
	_ = godotenv.Load()

	originURL := os.Getenv("OTD_ORIGIN_URL")
	if originURL == "" {
		originURL = "https://github.com/gabezeck/otd-cli"
	}
	contactEmail := os.Getenv("OTD_CONTACT_EMAIL")
	if contactEmail == "" {
		contactEmail = "contact@example.com"
	}

	url := "https://en.wikipedia.org/wiki/Wikipedia:On_this_day/Today"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Wikipedia requires a User-Agent header
	req.Header.Set("User-Agent", fmt.Sprintf("otd-cli/1.0 (%s; %s)", originURL, contactEmail))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	data := &OTDData{}

	mainContent := doc.Find(".mw-parser-output")

	// 1. Get Date
	// The date is typically the bold link that matches the pattern "Month Day".
	// We iterate through bold links to find one that looks like a date (not starting with "Wikipedia:")
	mainContent.Find("p b a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
		if !strings.HasPrefix(text, "Wikipedia:") {
			data.Date = text
			return false // break
		}
		return true // continue
	})

	if data.Date == "" {
		data.Date = "Today"
	}

	// 2. Get Events
	// The events are in a <ul>. On the 'Today' page, it is typically the first <ul>
	// that is a direct child of .mw-parser-output, or at least the first significant one.
	// We can iterate over all ULs and check the content format to be sure,
	// or just trust it's the first one as per standard layout.
	// We'll try the first one.
	eventsList := mainContent.Find("ul").First()
	eventsList.Find("li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		// Separators: en-dash is standard, but check others
		var parts []string
		if strings.Contains(text, "–") {
			parts = strings.SplitN(text, "–", 2)
		} else if strings.Contains(text, "-") {
			parts = strings.SplitN(text, "-", 2)
		}

		if len(parts) >= 2 {
			data.Events = append(data.Events, Event{
				Year: strings.TrimSpace(parts[0]),
				Text: strings.TrimSpace(parts[1]),
			})
		} else {
			// Keep it even if parse fails, just put whole text in Text
			if strings.TrimSpace(text) != "" {
				data.Events = append(data.Events, Event{
					Year: "",
					Text: strings.TrimSpace(text),
				})
			}
		}
	})

	// 3. Get Birthdays
	// These are in a div class="hlist".
	// We look for list items that contain "b." (born) or "d." (died).
	mainContent.Find(".hlist ul li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "(b.") || strings.Contains(text, "(d.") {
			name := s.Find("b").Text()
			if name == "" {
				// Sometimes the name isn't bold, take text before '('
				if idx := strings.Index(text, "("); idx != -1 {
					name = strings.TrimSpace(text[:idx])
				} else {
					name = text
				}
			}

			// Extract year info (the part in parens)
			yearInfo := ""
			if idx := strings.LastIndex(text, "("); idx != -1 {
				yearInfo = text[idx:]
				yearInfo = strings.Trim(yearInfo, "()")
			}

			data.Birthdays = append(data.Birthdays, Birthday{
				Name: name,
				Year: yearInfo,
			})
		}
	})

	return data, nil
}
