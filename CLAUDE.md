# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based token service that provides secure token creation, encryption, and retrieval capabilities. The service uses AES-GCM encryption and SHA512-256 hashing for token generation.

## Architecture

- **Entry Point**: `cmd/service/main.go` - HTTP server listening on port 8000
- **API Layer**: `api/` - HTTP handlers using Huma framework with Gorilla Mux router
  - `api/api.go` - Base handler with in-memory store
  - `api/routes.go` - Route registration and Huma API setup
  - `api/tokens.go` - Token-specific endpoints (create, get encrypted, get decrypted)
- **Models**: `models/` - Data structures and business logic
  - `models/models.go` - BaseModel with common fields (ID, timestamps)
  - `models/token.go` - Token model with encryption/decryption and tokenization logic
- **Infrastructure**: `intrastructure/` - AWS CDK deployment configuration (TypeScript)

## Key Components

### Token Flow
1. **Creation**: POST `/token` - Creates token from payload, encrypts it, generates SHA512-256 hash as token ID
2. **Retrieval**: GET `/token/{token}` - Returns encrypted token without decrypted payload
3. **Decryption**: GET `/token/{token}/decrypt` - Returns token with decrypted payload

### Data Storage
- Currently uses in-memory map storage (`BaseHandler.Store`)
- AWS infrastructure includes DynamoDB table and KMS key (not yet integrated)

## Common Commands

### Development
```bash
# Run the service
go run cmd/service/main.go

# Build the service
go build -o tokenize cmd/service/main.go

# Run tests
go test ./...

# Run tests for specific package
go test ./models

# Format code
go fmt ./...

# Vet code for common issues
go vet ./...
```

### Infrastructure (CDK)
```bash
cd intrastructure
npm run build
npm run test
npx cdk deploy
npx cdk diff
npx cdk synth
```

## Security Notes

- Hardcoded encryption key in `models/token.go:11` - should use environment variable or secure key management
- AES-GCM encryption with zero nonce (security concern in `models/token.go:41`)
- In-memory storage means tokens are lost on service restart

## Testing
- uses github.com/stretchr/testify/assert for assertions
- looking for 80% code coverage

## Dependencies

- **Huma v2**: Modern HTTP API framework for Go
- **Gorilla Mux**: HTTP router
- **Google UUID**: UUID generation (v7)
- **AWS CDK**: Infrastructure as code (TypeScript)