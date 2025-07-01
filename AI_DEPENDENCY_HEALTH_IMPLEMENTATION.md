# AI-Powered Dependency Health Agent Implementation

## Overview

Successfully implemented an AI-powered Dependency Health Agent for SBOM Sentinel that uses a local LLM via Ollama to assess the health of each component in an SBOM.

## Implementation Details

### 1. Dependency Health Agent (`internal/analysis/dependency_health.go`)

Created a new `DependencyHealthAgent` struct that implements the `AnalysisAgent` interface with the following features:

- **AI Integration**: Uses Ollama API (`http://localhost:11434/api/generate`) with the `llama3` model
- **Smart Prompting**: Generates specific prompts asking the LLM to assess project health and maintenance status
- **Risk Detection**: Analyzes LLM responses for risk indicators like "unmaintained", "deprecated", "risky", etc.
- **Error Handling**: Gracefully handles API failures and continues with other components
- **Structured Responses**: Parses JSON responses from Ollama API

#### Key Methods:
- `NewDependencyHealthAgent()`: Creates a new agent instance
- `Name()`: Returns "Dependency Health Agent"
- `Analyze()`: Iterates through SBOM components and queries LLM for each
- `generatePrompt()`: Creates specific prompts for component health assessment
- `queryOllama()`: Handles HTTP communication with Ollama API
- `indicatesRisk()`: Analyzes LLM responses for risk keywords

### 2. CLI Integration (`cmd/sentinel-cli/cmd/analyze.go`)

Enhanced the analyze command with:

- **New Flag**: `--enable-ai-health-check` to control AI analysis execution
- **Multi-Agent Support**: Runs both License Agent and Dependency Health Agent
- **Unified Results**: Combines results from all agents for display
- **User Feedback**: Provides hints when AI analysis is not enabled
- **Updated Documentation**: Enhanced command description and help text

#### Key Changes:
- Added flag registration for `--enable-ai-health-check`
- Modified analysis workflow to support multiple agents
- Updated result display to show combined findings
- Added core package import for `AnalysisResult` type

### 3. API Communication

The agent communicates with Ollama using:

```json
{
  "model": "llama3",
  "prompt": "Analyze the project health of...",
  "stream": false
}
```

And processes responses containing:
- AI-generated health assessment
- Metadata about the analysis (timing, token counts, etc.)

### 4. Risk Assessment

The agent identifies potential risks by checking LLM responses for keywords such as:
- "unmaintained", "deprecated", "risky", "outdated"
- "abandoned", "inactive", "archived", "obsolete"
- "discontinued", "end of life", "unsupported"
- "vulnerable", "security issues", "not recommended"

## Usage

```bash
# Standard analysis (license only)
./bin/sentinel-cli analyze sbom.json

# With AI-powered dependency health analysis
./bin/sentinel-cli analyze sbom.json --enable-ai-health-check

# Verbose output with AI analysis
./bin/sentinel-cli analyze sbom.json --enable-ai-health-check --verbose
```

## Requirements

- **Ollama**: Must be running locally on `localhost:11434`
- **Model**: Requires `llama3` model to be available
- **Network**: HTTP access to local Ollama instance

## Error Handling

- Gracefully handles Ollama API failures
- Continues analysis even if individual components fail
- Provides warning messages for failed analyses
- Maintains functionality when AI analysis is disabled

## Architecture Benefits

- **Modular Design**: Follows existing hexagonal architecture pattern
- **Interface Compliance**: Implements `AnalysisAgent` interface
- **Separation of Concerns**: AI logic isolated in dedicated agent
- **Extensible**: Easy to add more AI-powered analysis agents
- **Optional Feature**: Can be disabled when Ollama is not available

## Future Enhancements

Potential improvements for future iterations:
- Configurable LLM model selection
- Custom Ollama endpoint configuration
- Batch processing for improved performance
- Caching of AI responses
- More sophisticated risk scoring
- Support for additional LLM providers