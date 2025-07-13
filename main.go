package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Link represents a shortcut and its destination URL
type Link struct {
	Shortcut string `json:"shortcut"`
	URL      string `json:"url"`
}

// LinkStore manages the storage and retrieval of links
type LinkStore struct {
	links    map[string]string
	filePath string
}

// Server handles HTTP requests
type Server struct {
	store *LinkStore
}

// Load reads links from the JSON file
func (ls *LinkStore) Load() error {
	// Ensure directory exists
	dir := filepath.Dir(ls.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(ls.filePath); os.IsNotExist(err) {
		// File doesn't exist, start with empty map
		return nil
	}

	// Read the file
	data, err := os.ReadFile(ls.filePath)
	if err != nil {
		return err
	}

	// Parse JSON
	var links []Link
	if err := json.Unmarshal(data, &links); err != nil {
		return err
	}

	// Convert to map
	for _, link := range links {
		ls.links[link.Shortcut] = link.URL
	}

	return nil
}

// Save writes links to the JSON file
func (ls *LinkStore) Save() error {
	// Convert map to slice
	var links []Link
	for shortcut, url := range ls.links {
		links = append(links, Link{
			Shortcut: shortcut,
			URL:      url,
		})
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(links, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(ls.filePath, data, 0644)
}

// Add creates a new link
func (ls *LinkStore) Add(shortcut, url string) error {
	ls.links[shortcut] = url
	return ls.Save()
}

// Get retrieves a URL by shortcut
func (ls *LinkStore) Get(shortcut string) (string, bool) {
	url, exists := ls.links[shortcut]
	return url, exists
}

// GetAll returns all links
func (ls *LinkStore) GetAll() map[string]string {
	result := make(map[string]string)
	for k, v := range ls.links {
		result[k] = v
	}
	return result
}

// handleHome handles the homepage and redirect requests
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	// If path is empty, show homepage
	if path == "" {
		s.showHomepage(w, r)
		return
	}

	// Try to redirect to the URL for this shortcut
	if url, exists := s.store.Get(path); exists {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	// Shortcut not found, redirect to homepage
	http.Redirect(w, r, "/", http.StatusFound)
}

// handleAdd handles form submissions to add new links
func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	shortcut := strings.TrimSpace(r.FormValue("shortcut"))
	url := strings.TrimSpace(r.FormValue("url"))

	// Basic validation
	if shortcut == "" || url == "" {
		http.Error(w, "Shortcut and URL are required", http.StatusBadRequest)
		return
	}

	// Add http:// if no protocol specified
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// Save the new link
	if err := s.store.Add(shortcut, url); err != nil {
		http.Error(w, "Failed to save link", http.StatusInternalServerError)
		return
	}

	// Redirect back to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// showHomepage renders the HTML homepage
func (s *Server) showHomepage(w http.ResponseWriter, r *http.Request) {
	const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Links</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            background-color: #f8f9fa;
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 2rem;
        }
        .form-group {
            margin-bottom: 1rem;
        }
        label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #555;
        }
        input[type="text"], input[type="url"] {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 1rem;
            box-sizing: border-box;
        }
        button {
            background-color: #007bff;
            color: white;
            padding: 0.75rem 2rem;
            border: none;
            border-radius: 4px;
            font-size: 1rem;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        button:hover {
            background-color: #0056b3;
        }
        .links-section {
            margin-top: 3rem;
        }
        .links-list {
            background: #f8f9fa;
            border-radius: 4px;
            padding: 1rem;
        }
        .link-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 0.75rem;
            margin: 0.5rem 0;
            background: white;
            border-radius: 4px;
            border: 1px solid #e9ecef;
        }
        .shortcut {
            font-weight: 600;
            color: #007bff;
            font-family: monospace;
        }
        .url {
            color: #666;
            word-break: break-all;
        }
        .empty-state {
            text-align: center;
            color: #666;
            font-style: italic;
            padding: 2rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔗 Go Links</h1>
        
        <form action="/add" method="post">
            <div class="form-group">
                <label for="shortcut">Shortcut:</label>
                <input type="text" id="shortcut" name="shortcut" placeholder="e.g., gh" required>
            </div>
            <div class="form-group">
                <label for="url">URL:</label>
                <input type="url" id="url" name="url" placeholder="e.g., https://github.com" required>
            </div>
            <button type="submit">Add Link</button>
        </form>

        <div class="links-section">
            <h2>Your Links</h2>
            <div class="links-list">
                {{if .Links}}
                    {{range $shortcut, $url := .Links}}
                    <div class="link-item">
                        <span class="shortcut">go/{{$shortcut}}</span>
                        <span class="url">→ {{$url}}</span>
                    </div>
                    {{end}}
                {{else}}
                    <div class="empty-state">
                        No links yet. Add your first one above!
                    </div>
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("homepage").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Links map[string]string
	}{
		Links: s.store.GetAll(),
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize the link store
	store := &LinkStore{
		links:    make(map[string]string),
		filePath: "/app/data/links.json",
	}

	// Load existing links from file
	if err := store.Load(); err != nil {
		log.Printf("Warning: Could not load links file: %v", err)
	}

	// Initialize the server
	server := &Server{store: store}

	// Set up routes
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/add", server.handleAdd)

	// Start the server
	fmt.Println("Go Links server starting on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
