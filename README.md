# SBOM Sentinel

![Development Status](https://img.shields.io/badge/status-Alpha-red?style=flat-square)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8)](https://golang.org/)
[![SQLite](https://img.shields.io/badge/database-SQLite-003B57)](https://sqlite.org/)

**SBOM Sentinel** is a Software Bill of Materials (SBOM) analysis platform that provides deep, contextual intelligence on software supply chain risks through AI-powered analysis and traditional security scanning.

## Why SBOM Sentinel?

SBOM Sentinel goes beyond CVE databases by leveraging local AI models to assess dependency health, maintenance status, and emerging risks that haven't yet been cataloged in vulnerability databases. It provides both immediate actionable insights and strategic intelligence about your software supply chain.

## 🚀 Core Features

- **📄 CycloneDX SBOM Parsing** - Complete support for industry-standard SBOM format
- **⚖️ License Compliance Analysis** - Automated detection of high-risk copyleft licenses
- **🤖 AI-Powered Dependency Health Checks** - Intelligent assessment using local Ollama LLM
- **🔍 Proactive Vulnerability Discovery** - RAG-powered detection of pre-CVE threats from security intelligence
- **💾 SQLite-based Persistence** - Efficient storage and retrieval of SBOM documents
- **🔄 Dual Interface** - Both command-line tool and REST API server
- **🏗️ Hexagonal Architecture** - Clean, testable, and extensible codebase design
- **📊 Comprehensive Analysis Results** - Detailed findings with severity classification

## 🎯 Quick Start Guide

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/hueyexe/SBOM-Sentinel.git
   cd SBOM-Sentinel
   ```

2. **Build the applications:**
   ```bash
   # Build both CLI and server
   go build -o bin/sentinel-cli ./cmd/sentinel-cli/
   go build -o bin/sentinel-server ./cmd/sentinel-server/
   ```

3. **Verify installation:**
   ```bash
   ./bin/sentinel-cli --help
   ./bin/sentinel-server --help
   ```

### CLI Usage

#### Basic SBOM Analysis
```bash
# Analyze an SBOM file with license compliance checking
./bin/sentinel-cli analyze your-sbom.json

# Show detailed component information
./bin/sentinel-cli analyze your-sbom.json --verbose

# Generate summary report only
./bin/sentinel-cli analyze your-sbom.json --summary
```

#### AI-Powered Analysis
```bash
# Enable AI-powered dependency health analysis (requires Ollama)
./bin/sentinel-cli analyze your-sbom.json --enable-ai-health-check

# Enable proactive vulnerability discovery using RAG (requires Ollama)
./bin/sentinel-cli analyze your-sbom.json --enable-proactive-scan

# Combine all AI features with verbose output
./bin/sentinel-cli analyze your-sbom.json --enable-ai-health-check --enable-proactive-scan --verbose
```

**Example Output:**
```
✅ Successfully parsed SBOM: MyApplication v1.0.0
📦 Found 42 components

🔬 Analysis Results:
   Found 2 issues:

   1. 🟡 [Medium] License Agent
      Component 'copyleft-library' (v2.1.0) uses high-risk copyleft license 'GPL-3.0-only'

   2. 🟡 [Medium] Dependency Health Agent
      The project 'abandoned-utility' appears to be unmaintained and deprecated

   3. 🟡 [Medium] Proactive Vulnerability Agent
      Component 'example-lib' may be vulnerable to deserialization issues based on security intelligence
```

### API Usage

#### 1. Start the Server
```bash
# Start with default settings (port 8080, local database)
./bin/sentinel-server

# Or configure with environment variables
DATABASE_PATH=/path/to/sentinel.db PORT=9000 ./bin/sentinel-server
```

#### 2. Submit an SBOM
```bash
# Upload an SBOM file for storage
curl -X POST \
  -F "sbom=@your-sbom.json" \
  http://localhost:8080/api/v1/sboms

# Example response:
# {
#   "id": "urn:uuid:12345678-1234-1234-1234-123456789012",
#   "message": "SBOM submitted successfully"
# }
```

#### 3. Analyze the Stored SBOM
```bash
# Run basic analysis (license compliance only)
curl -X POST \
  "http://localhost:8080/api/v1/sboms/urn:uuid:12345678-1234-1234-1234-123456789012/analyze"

# Run analysis with AI health checks
curl -X POST \
  "http://localhost:8080/api/v1/sboms/urn:uuid:12345678-1234-1234-1234-123456789012/analyze?enable-ai-health-check=true"

# Run analysis with proactive vulnerability discovery
curl -X POST \
  "http://localhost:8080/api/v1/sboms/urn:uuid:12345678-1234-1234-1234-123456789012/analyze?enable-proactive-scan=true"

# Run comprehensive analysis with all AI features
curl -X POST \
  "http://localhost:8080/api/v1/sboms/urn:uuid:12345678-1234-1234-1234-123456789012/analyze?enable-ai-health-check=true&enable-proactive-scan=true"
```

**Example Analysis Response:**
```json
{
  "sbom_id": "urn:uuid:12345678-1234-1234-1234-123456789012",
  "results": [
    {
      "agent_name": "License Agent",
      "finding": "Component 'example-lib' (v1.0.0) uses high-risk copyleft license 'GPL-3.0-only'",
      "severity": "High"
    },
    {
      "agent_name": "Proactive Vulnerability Agent",
      "finding": "Component 'data-processor' may be vulnerable to memory leaks based on security discussions",
      "severity": "Medium"
    }
  ],
      "summary": {
      "total_findings": 2,
      "findings_by_severity": {
        "High": 1,
        "Medium": 1
      },
      "agents_run": ["License Agent", "Proactive Vulnerability Agent"]
    }
}
```

#### 4. Retrieve Stored SBOMs
```bash
# Get SBOM by ID
curl "http://localhost:8080/api/v1/sboms/get?id=urn:uuid:12345678-1234-1234-1234-123456789012"

# Health check
curl http://localhost:8080/health
```

## 🧠 AI-Powered Analysis

SBOM Sentinel integrates with [Ollama](https://ollama.ai/) to provide intelligent dependency health assessments:

### Prerequisites
1. **Install Ollama:** Follow the [official installation guide](https://ollama.ai/)
2. **Download a model:** `ollama pull llama3`
3. **Start Ollama:** `ollama serve`

### How It Works

**Dependency Health Analysis:**
The AI agent analyzes each component by:
- Querying the LLM about project health and maintenance status
- Detecting risk indicators like "unmaintained", "deprecated", "abandoned"
- Providing contextual insights beyond traditional vulnerability databases
- Flagging components that may pose supply chain risks

**Proactive Vulnerability Discovery:**
The RAG-powered agent provides early threat detection by:
- Harvesting security intelligence from discussions, forums, and research
- Creating vector embeddings of security documents using local AI
- Performing similarity searches against component names and versions
- Using LLM analysis to identify potential vulnerabilities before CVE publication
- Discovering emerging threats from unstructured security data sources

## 📋 Supported SBOM Formats

Currently supported:
- **CycloneDX JSON** (v1.4+)

Planned support:
- SPDX JSON/YAML
- SWID Tags
- Custom JSON formats

## 🏗️ Architecture

SBOM Sentinel follows a **hexagonal (ports and adapters) architecture**:

```
┌─────────────────────────────────────────────────────────┐
│                    Transport Layer                      │
│  ┌─────────────┐              ┌─────────────────────┐  │
│  │   CLI Tool  │              │    REST API Server  │  │
│  └─────────────┘              └─────────────────────┘  │
├─────────────────────────────────────────────────────────┤
│                     Core Domain                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │   Models    │  │  Analysis   │  │   Ingestion     │ │
│  │   (SBOM)    │  │   Agents    │  │    Parsers      │ │
│  └─────────────┘  └─────────────┘  └─────────────────┘ │
├─────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                   │
│  ┌─────────────┐              ┌─────────────────────┐  │
│  │   SQLite    │              │     Ollama API      │  │
│  │  Database   │              │   (AI Integration)  │  │
│  └─────────────┘              └─────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_PATH` | SQLite database file path | `./sentinel.db` |

### CLI Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Enable detailed output |
| `--summary` | Show summary only |
| `--format` | SBOM format (auto, cyclonedx) |
| `--enable-ai-health-check` | Enable AI health analysis |
| `--enable-proactive-scan` | Enable RAG-based vulnerability discovery |

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🚧 Roadmap

- [ ] **SPDX format support**
- [ ] **PostgreSQL backend option**
- [x] **Web dashboard interface**
- [ ] **Vulnerability database integration**
- [ ] **Custom analysis rule engine**
- [x] **Batch processing capabilities**
- [ ] **Export to various report formats**
- [ ] **Integration with CI/CD pipelines**