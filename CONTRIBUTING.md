# Contributing to SBOM Sentinel

Thank you for your interest in contributing to SBOM Sentinel! This document provides guidelines for contributing to the project.

## 🚀 Getting Started

### Prerequisites

- Go 1.24+ for backend development
- Node.js 18+ and npm for frontend development
- Git for version control

### Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/your-username/SBOM-Sentinel.git
   cd SBOM-Sentinel
   ```

2. **Set up the backend:**
   ```bash
   go mod download
   go build ./...
   ```

3. **Set up the frontend:**
   ```bash
   cd web
   npm install
   ```

## 🏗️ Project Structure

```
SBOM-Sentinel/
├── cmd/                    # Command-line applications
│   ├── sentinel-cli/      # CLI tool
│   └── sentinel-server/   # REST API server
├── internal/              # Core application logic
│   ├── analysis/          # Analysis agents
│   ├── core/              # Domain models
│   ├── ingestion/         # SBOM parsing
│   ├── platform/          # Infrastructure
│   └── transport/         # HTTP handlers
├── web/                   # React frontend
│   ├── src/
│   │   ├── components/    # Reusable UI components
│   │   ├── pages/         # Route components
│   │   ├── services/      # API integration
│   │   └── types/         # TypeScript definitions
│   └── tests/             # Frontend tests
└── docs/                  # Documentation
```

## 🧪 Testing

### Backend Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./internal/analysis/
```

### Frontend Testing
```bash
cd web

# Run unit tests
npm test

# Run E2E tests
npm run test:e2e

# Run tests with coverage
npm run test:coverage
```

## 📝 Code Style

### Go (Backend)
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Write tests for new functionality
- Use meaningful variable and function names
- Add comments for exported functions

### TypeScript/React (Frontend)
- Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- Use TypeScript for all new code
- Write functional components with hooks
- Use meaningful component and variable names
- Add JSDoc comments for complex functions

## 🔧 Development Workflow

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Write code following the style guidelines
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes:**
   ```bash
   # Backend
   go test ./...
   go build ./...

   # Frontend
   cd web
   npm test
   npm run build
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and create a pull request:**
   ```bash
   git push origin feature/your-feature-name
   ```

## 📋 Pull Request Guidelines

### Before submitting a PR:

1. **Ensure code quality:**
   - All tests pass
   - Code follows style guidelines
   - No linting errors
   - Documentation is updated

2. **Write a clear description:**
   - What the PR does
   - Why the changes are needed
   - How to test the changes
   - Any breaking changes

3. **Include relevant information:**
   - Screenshots for UI changes
   - Test results
   - Performance impact (if applicable)

### PR Template:
```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment details:**
   - Operating system
   - Go version
   - Node.js version
   - Browser (for frontend issues)

2. **Steps to reproduce:**
   - Clear, step-by-step instructions
   - Sample SBOM file (if applicable)
   - Expected vs actual behavior

3. **Additional information:**
   - Error messages
   - Logs
   - Screenshots

## 💡 Feature Requests

When requesting features, please include:

1. **Problem description:**
   - What problem does this solve?
   - Who would benefit from this feature?

2. **Proposed solution:**
   - How should the feature work?
   - Any design considerations?

3. **Additional context:**
   - Related issues or discussions
   - Similar features in other tools

## 📚 Documentation

### Contributing to Documentation

- Keep documentation up to date with code changes
- Use clear, concise language
- Include code examples where helpful
- Update README.md for significant changes

### Documentation Structure

- `README.md` - Project overview and quick start
- `docs/` - Detailed documentation
- Code comments - Inline documentation

## 🤝 Community Guidelines

### Be Respectful
- Be kind and respectful to all contributors
- Welcome newcomers and help them get started
- Provide constructive feedback

### Be Collaborative
- Share knowledge and help others
- Review PRs promptly and constructively
- Participate in discussions and decisions

### Be Professional
- Follow the project's code of conduct
- Maintain high code quality standards
- Communicate clearly and professionally

## 📄 License

By contributing to SBOM Sentinel, you agree that your contributions will be licensed under the MIT License.

## 🆘 Getting Help

- **GitHub Issues:** For bug reports and feature requests
- **GitHub Discussions:** For questions and general discussion
- **Documentation:** Check the README and docs folder first

---

Thank you for contributing to SBOM Sentinel! 🚀