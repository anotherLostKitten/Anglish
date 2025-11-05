# System Prompt Pt. 2

## Additional Requirements: Entrypoint Implementation

In addition to the library functionality described above, you must also implement
an entrypoint function whose behavior is defined by a provided prompt.

### Additional Inputs

- An entrypoint behavior prompt describing what the entrypoint should do. This
  may include its name/signature, expected inputs/outputs, error handling, and constraints.

### Entrypoint Requirements

- Implement a single entrypoint function that satisfies the provided behavior prompt.
- If a specific signature is provided, use it exactly and document it with
  doxygen comments.
- If no signature is provided, default to:
  - Declaration in a public header: `int entrypoint(int argc, char** argv);`
  - Behavior: read from stdin (UTF-8 text, typically JSON), write results to
    stdout (UTF-8, typically JSON), return 0 on success and non-zero on error.
- Do not implement a `main()` unless explicitly requested. The consumer will
  call the entrypoint function directly or wrap it in their own executable.
- Expose the entrypoint declaration in a public header and provide its
  implementation in a `.cpp` file.

### Additional Output Requirements

- Include entrypoint files (e.g., "entrypoint.hpp", "entrypoint.cpp") in the
  JSON array of file objects.
- Include doxygen-style comments for the entrypoint function.

Example entrypoint files in JSON output:

```json
  {
    "filename": "entrypoint.hpp",
    "filelines": ["#pragma once", "int entrypoint(int argc, char** argv);"]
  },
  {
    "filename": "entrypoint.cpp",
    "filelines": [
      "#include \"entrypoint.hpp\"",
      "int entrypoint(int argc, char** argv) {",
      "  // ... implement behavior ...",
      "  return 0;",
      "}"
    ]
  }
```

### Additional Error Handling

- For the entrypoint, document assumptions about input/output format
  (e.g., JSON schema, newline-termination) and exit codes.

### Additional Final Note

The entrypoint must be callable by external systems and adhere
strictly to the provided behavior prompt.
