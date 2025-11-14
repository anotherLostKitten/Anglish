package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anotherLostKitten/Anglish/internal/verify"
)

func TestValidateGoResponseHelloWorld(t *testing.T) {
	helloWorldJSON := []byte(`[
  {
    "filename": "cmd/helloworld/main.go",
    "filelines": [
      "package main",
      "",
      "import \"fmt\"",
      "",
      "func main() {",
      "    fmt.Println(\"hello, world\")",
      "}"
    ]
  }
]`)

	require.NoError(t, verify.ValidateGoResponse(helloWorldJSON))
}

func TestValidateGoResponseSyntaxError(t *testing.T) {
	badProgramJSON := []byte(`[
  {
    "filename": "cmd/bad/main.go",
    "filelines": [
      "package main",
      "",
      "func main() {",
      "    fmt.Println(\"oops\"",
      "}"
    ]
  }
]`)

	require.Error(t, verify.ValidateGoResponse(badProgramJSON))
}
