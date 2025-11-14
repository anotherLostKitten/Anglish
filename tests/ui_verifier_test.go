package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anotherLostKitten/Anglish/internal/verify"
)

func TestValidateJSResponseHelloWorld(t *testing.T) {
	helloWorldJS := []byte(`[
  {
    "filename": "web/app.js",
    "filelines": [
      "function main() {",
      "  console.log('Hello, world!');",
      "}",
      "",
      "main();"
    ]
  }
]`)

	require.NoError(t, verify.ValidateJSResponse(helloWorldJS))
}

func TestValidateHTMLResponseValid(t *testing.T) {
	validHTML := []byte(`[
  {
    "filename": "web/index.html",
    "filelines": [
      "<!doctype html>",
      "<html lang=\"en\">",
      "<head>",
      "  <meta charset=\"utf-8\">",
      "  <title>Hello World</title>",
      "</head>",
      "<body>",
      "  <main>",
      "    <h1>Hello, World!</h1>",
      "    <p>This is a minimal document.</p>",
      "  </main>",
      "</body>",
      "</html>"
    ]
  }
]`)

	require.NoError(t, verify.ValidateHTMLResponse(validHTML))
}

func TestValidateHTMLResponseInvalid(t *testing.T) {
	invalidHTML := []byte(`[
  {
    "filename": "web/index.html",
    "filelines": []
  }
]`)

	require.Error(t, verify.ValidateHTMLResponse(invalidHTML))
}

func TestValidateCSSResponseValid(t *testing.T) {
	validCSS := []byte(`[
  {
    "filename": "web/styles.css",
    "filelines": [
      ":root {",
      "  --color-text: #111;",
      "}",
      "",
      "main {",
      "  min-height: 100vh;",
      "}"
    ]
  }
]`)

	require.NoError(t, verify.ValidateCSSResponse(validCSS))
}

func TestValidateCSSResponseInvalid(t *testing.T) {
	invalidCSS := []byte(`[
  {
    "filename": "web/styles.css",
    "filelines": []
  }
]`)

	require.Error(t, verify.ValidateCSSResponse(invalidCSS))
}

func TestValidateJSResponseSyntaxError(t *testing.T) {
	payload := []byte(`[
  {
    "filename": "ui/app.js",
    "filelines": [
      "function broken() {",
      "  console.log('missing bracket'",
      "",
      "broken();"
    ]
  }
]`)

	err := verify.ValidateJSResponse(payload)
	require.Error(t, err)
}
