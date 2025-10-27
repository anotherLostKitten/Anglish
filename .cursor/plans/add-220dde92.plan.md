<!-- 220dde92-6a0a-4319-b903-20371a7adf51 b051cbce-9189-428c-9a03-0e77849f0b02 -->
# Add Agentic Space System Prompt

## What we will do

- Create a new prompt at `internal/llm/prompts/agentic_space_agent_system.md` modeled after `func_space_agent_system.md`.
- Keep all C++ library generation guidance, but adapt for agentic orchestration use.
- Add explicit output requirements for a separate JSON file containing OpenAI Tools API function schemas (type:function) for every callable function.

## Key content of the new prompt

- Purpose: Generate C++ agentic-space libraries (functions/classes/data) designed for LLM invocation.
- Inputs always provided: goals, required elements, any APIs.
- Output Requirements:
- C++ code: headers/implementations, namespaces/classes, doxygen comments, minimal usage examples.
- JSON schema file: `agentic_tools.json` containing an array of tools entries, one per function, using OpenAI Tools API format:
- Each entry: `{ "type": "function", "function": { "name", "description", "parameters" } }`
- `parameters` is a strict JSON Schema object: `{ "$schema": "http://json-schema.org/draft-07/schema#", "type":"object", "properties":{}, "required":[], "additionalProperties": false }`.
- Agentic-specific design rules:
- Make functions LLM-friendly: clear names, single responsibility, stable signatures, deterministic where possible.
- Avoid overloading; prefer unique names per behavior.
- Prefer pure functions; isolate side effects; validate inputs; precise errors.
- No hidden global state; thread-safe where applicable.
- Mapping rules C++ → JSON Schema:
- int/long → integer; double/float → number; bool → boolean; std::string → string
- enum → string with `enum` values; std::vector<T> → array with `items`
- optional<T> → property optional; struct → object with nested properties
- Constrain with `minimum`, `maximum`, `pattern`, `minItems`, `maxItems` when applicable.
- Consistency Rules: keep existing style (snake_case, PascalCase types, RAII, const correctness, exceptions, C++17+).
- Error Handling: validate inputs, state assumptions in comments, resolve conflicts logically.
- Output formatting guidance: emit code files and a single JSON file; names in JSON must exactly match C++ function names; keep `additionalProperties:false`.

## Files to add

- `internal/llm/prompts/agentic_space_agent_system.md`: Complete prompt content described above.

### To-dos

- [ ] Add `internal/llm/prompts/agentic_space_agent_system.md` with agentic-space prompt and Tools API JSON rules
- [ ] Define precise C++→JSON Schema mapping and constraints in the prompt