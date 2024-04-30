package handlers

import (
	"github.com/russross/blackfriday"
	"net/http"
	"os"
)

// DefaultHandler
// Default handler for the server. Redirects to the web page.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// If the request is for the root path, redirect to the web page
	if r.URL.Path == "/" {
		// Read the contents of the README.md file
		readme, err := os.ReadFile("README.md")
		if err != nil {
			http.Error(w, "Failed to read README.md", http.StatusInternalServerError)
			return
		}

		html := blackfriday.MarkdownCommon(readme)

		// Set the Content-Type header to indicate that this is HTML
		w.Header().Set("Content-Type", "text/html")

		// Custom styles for the HTML to make it look better
		customStyles := `
            <style>
                pre {background-color: #f4f4f4;}
				body {padding: 10px;}
            </style>
        `

		// Append the custom styles to the HTML
		htmlWithStyles := append([]byte(customStyles), html...)

		// Write the HTML response
		_, err = w.Write(htmlWithStyles)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
