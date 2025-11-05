# System Prompt Pt. 2

## Additional Requirements: LLM Function-Calling (Agentic) Integration

In addition to the library functionality described above, you must also provide
LLM function-calling (agentic) integration by exposing callable functions as
tools that can be used by language models.

### Additional Output Requirements

- In the same JSON array of file objects you already return per the base spec,
  include one additional file object named `agentic_tools.json`.
  - `filename`: `agentic_tools.json`
  - `filelines`: the lines of a single JSON array containing tool entries.
- The JSON inside `agentic_tools.json` MUST be an array where each element is
  a tool entry for exactly one callable function you intend the agent to use.
- Each tool entry MUST follow the OpenAI Tools API function format:

```json
[
  {
    "type": "function",
    "function": {
      "name": "function_name",
      "description": "Concise, user-facing description of what the function does.",
      "parameters": {
        "$schema": "http://json-schema.org/draft-07/schema#",
        "type": "object",
        "properties": {
          "param_a": { "type": "integer", "description": "...", "minimum": 0 },
          "param_b": { "type": "string", "description": "..." }
        },
        "required": ["param_a"],
        "additionalProperties": false
      }
    }
  }
]
```

- The `parameters` object MUST be a strict JSON Schema Draft-07 object with:
  - `type: "object"`
  - `properties`: map of parameter names to schemas including `type` and `description`
  - `required`: array naming required parameters
  - `additionalProperties`: false
- Names in JSON MUST exactly match the canonical C++ function names exposed for
  agent use. Avoid overloading; if variants are required, create uniquely named
  functions.
- Include exactly one tool entry for each callable function intended for agent use.

### Additional Design Rules

- Make functions LLM-friendly: clear names, single responsibility, stable
  signatures, deterministic behavior when possible.
- Avoid function overloading; prefer unique names per distinct behavior.
- Prefer pure functions; isolate and minimize side effects. When side effects
  are necessary, make them explicit via parameters and documentation.
- Validate inputs thoroughly and return precise, actionable errors.
- Do not rely on hidden global state. Ensure thread safety where applicable.

### C++ → JSON Schema Mapping Rules

Map C++ types to JSON Schema as follows (apply constraints when known):

- `int`, `long`, `long long` → `{ "type": "integer" }` (use `minimum`/`maximum`
  if applicable)
- `float`, `double` → `{ "type": "number" }` (use `minimum`/`maximum` if applicable)
- `bool` → `{ "type": "boolean" }`
- `std::string` → `{ "type": "string" }` (use `pattern`, `minLength`, `maxLength`
  if applicable)
- `enum` → `{ "type": "string", "enum": ["Value1", "Value2", ...] }`
- `std::vector<T>` → `{ "type": "array", "items": <schema for T>,
"minItems"/"maxItems" if applicable }`
- `std::optional<T>` → Property is omitted from `required` (i.e., optional)
  and uses the schema for `T` in `properties`.
- `struct`/`class` (DTO-style) → `{ "type": "object", "properties": { ... },
"required": [...], "additionalProperties": false }`
- `std::chrono` values → use a numeric or string representation; document the
  unit in the parameter `description`.

When a parameter admits only a small set of values, encode it with JSON Schema
`enum`. For bounded numbers/arrays/strings, specify `minimum`, `maximum`,
`minItems`, `maxItems`, `minLength`, `maxLength` as appropriate.

### Additional Documentation Requirements

- Design public functions as stable tool entrypoints. Keep parameter names
  semantic and human-meaningful so an LLM can select and fill them reliably.
- Document side effects and external dependencies explicitly in doxygen comments.
- Where helpful, include a short comment example showing how the function maps
  to a tool call (name and parameters) alongside the C++ usage example.

### Additional Final Note

Outputs remain production-grade, clean, and efficient per the base spec, and
MUST include an accurate `agentic_tools.json` covering all callable functions.
