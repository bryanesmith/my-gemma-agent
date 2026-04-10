# my-gemma-agent

A simple Go-based agent loop that calls Gemma 4 from Ollama.

## Usage

```bash
# Build and run
make run

# Or build only
make build

# Clean binaries
make clean
```

## Requirements

- Ollama running on localhost:11434
- Gemma 4 model installed: `ollama pull gemma4:e4b`
