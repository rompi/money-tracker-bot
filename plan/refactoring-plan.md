# Money Tracker Bot - Refactoring Plan

## Overview
This document outlines refactoring opportunities to improve code quality, maintainability, and reliability without changing existing features. The plan focuses on clean architecture principles, error handling improvements, and code organization enhancements.

## 🏗️ Architecture & Design Improvements

### 1. Configuration Management
**Current Issue**: Environment variables scattered across multiple files, hard to manage and test.

**Improvements**:
- Create a centralized `internal/config` package
- Define configuration struct with validation
- Support different environments (dev, staging, prod)
- Add configuration loading with defaults and validation

**Files to Refactor**:
- `cmd/telebot/main.go` (environment variable loading)
- `internal/adapters/telegram/handler.go` (spreadsheet ID access)
- `internal/service/transactions/handler.go` (spreadsheet ID access)

### 2. Error Handling & Logging
**Current Issue**: Inconsistent error handling with `log.Fatal` and `log.Panic` calls that crash the application.

**Improvements**:
- Replace `log.Fatal` and `log.Panic` with proper error returns
- Implement structured logging with levels (debug, info, warn, error)
- Create custom error types for different error categories
- Add error wrapping with context information
- Implement graceful degradation instead of crashes

**Critical Files**:
- `internal/adapters/google/spreadsheet/client.go:32,45,65,108`
- `internal/adapters/telegram/handler.go:32,60`
- `internal/adapters/gemini/gemini.go:43`

### 3. Context Management
**Current Issue**: Using `context.TODO()` instead of proper context propagation.

**Improvements**:
- Replace `context.TODO()` with proper context from request chain
- Add timeout and cancellation support for AI operations
- Implement context-aware operations throughout the stack

**Files to Fix**:
- `internal/adapters/telegram/handler.go:155,189`

## 🧪 Testing & Quality Improvements

### 4. Test Coverage Enhancement
**Current State**: Basic tests exist but coverage is incomplete.

**Improvements**:
- Add comprehensive unit tests for all packages
- Implement integration tests for critical workflows
- Add table-driven tests for data validation
- Create test helpers and fixtures
- Add benchmark tests for performance-critical operations

**Priority Areas**:
- Transaction validation logic
- AI response parsing
- Spreadsheet operations
- Error scenarios

### 5. Mock & Interface Improvements
**Current Issue**: Some concrete dependencies, limited mocking capabilities.

**Improvements**:
- Create interfaces for all external dependencies
- Improve mock implementations with better validation
- Add dependency injection container
- Separate test utilities into dedicated packages

## 🔧 Code Organization & Structure

### 6. Validation Layer
**Current Issue**: Input validation scattered across different layers.

**Improvements**:
- Create centralized validation package
- Add input sanitization for user messages
- Implement transaction data validation
- Add amount parsing and validation utilities
- Create validation middleware for API inputs

### 7. Constants & Magic Numbers
**Current Issue**: Hardcoded values throughout the codebase.

**Improvements**:
- Extract magic numbers to named constants
- Create constants file for API limits, timeouts, etc.
- Move business rules to configuration
- Add validation for predefined lists (categories, accounts)

**Examples**:
- Timeout values (60 seconds in telegram handler)
- Sheet ranges ("detailed!A:G", "summary!A2:F12")
- Timezone strings ("Asia/Bangkok")
- File paths and directories

### 8. File Management
**Current Issue**: File operations mixed with business logic.

**Improvements**:
- Create dedicated file management service
- Implement proper cleanup with defer statements
- Add file validation (size, type, security)
- Create temporary file management utilities
- Add file operation error handling

## 📊 Data & Persistence Improvements

### 9. Amount Handling
**Current Issue**: String-based amount handling without proper validation.

**Improvements**:
- Create dedicated Money/Amount type with proper decimal handling
- Add currency conversion utilities
- Implement amount formatting standardization
- Add validation for amount ranges and formats

### 10. Database Abstraction
**Current Issue**: Direct Google Sheets coupling without abstraction.

**Improvements**:
- Create repository pattern for data operations
- Add database interface to support multiple backends
- Implement transaction patterns for data consistency
- Add data migration utilities for schema changes

## 🔒 Security & Reliability Improvements

### 11. Input Sanitization
**Current Issue**: Limited input validation for user messages and file uploads.

**Improvements**:
- Add comprehensive input sanitization
- Implement file upload security checks
- Add rate limiting for API calls
- Create request validation middleware

### 12. Graceful Shutdown
**Current Issue**: No graceful shutdown handling.

**Improvements**:
- Implement signal handling for graceful shutdown
- Add cleanup for ongoing operations
- Ensure proper resource cleanup
- Add health check endpoints

## 🚀 Performance & Scalability

### 13. Concurrency & Performance
**Current Issue**: Sequential processing of operations.

**Improvements**:
- Add concurrent processing for independent operations
- Implement connection pooling for external services
- Add caching for frequently accessed data
- Optimize AI API call batching

### 14. Monitoring & Observability
**Current Issue**: Limited visibility into application behavior.

**Improvements**:
- Add metrics collection (processing time, success/failure rates)
- Implement distributed tracing
- Add application health monitoring
- Create performance dashboards

## 📁 Package Structure Reorganization

### 15. Clean Architecture Enforcement
**Current State**: Good foundation but some violations.

**Improvements**:
- Ensure strict dependency direction (domain <- service <- adapter)
- Move shared types to appropriate layers
- Create clear boundaries between layers
- Add architecture tests to prevent violations

### 16. Shared Utilities
**Current Issue**: Utilities mixed with business logic.

**Improvements**:
- Create `internal/utils` for pure functions
- Separate formatting utilities
- Create reusable validation helpers
- Add common middleware components

## 🔄 Refactoring Priority Matrix

### High Priority (Immediate)
1. **Error Handling** - Replace fatal errors with proper error returns
2. **Context Management** - Fix context.TODO() usage
3. **Configuration** - Centralize environment variable management
4. **Input Validation** - Add comprehensive validation layer

### Medium Priority (Next Sprint)
5. **Testing** - Increase test coverage significantly
6. **Constants** - Extract magic numbers and hardcoded values
7. **File Management** - Improve file operation handling
8. **Amount Handling** - Create proper money type

### Low Priority (Future)
9. **Monitoring** - Add observability features
10. **Performance** - Optimize for higher throughput
11. **Database Abstraction** - Support multiple backends
12. **Security** - Enhanced security measures

## 📋 Implementation Guidelines

### Refactoring Principles
1. **Backward Compatibility**: All refactoring must maintain existing API contracts
2. **Incremental Changes**: Implement changes in small, reviewable chunks
3. **Test-Driven**: Add tests before refactoring existing code
4. **Documentation**: Update documentation as changes are made
5. **Performance**: Ensure refactoring doesn't degrade performance

### Success Metrics
- Increased test coverage (target: >80%)
- Reduced cyclomatic complexity
- Improved error handling coverage
- Zero application crashes in production
- Faster development velocity for new features

---

## 📅 Phased Implementation Plan

This section breaks down the refactoring into 5 manageable phases, each building upon the previous phase to ensure stability and incremental improvement.

---

## 📋 Phase 1: Foundation & Stability (Week 1-2)
**Goal**: Fix critical stability issues and establish proper error handling foundation.

### 🎯 Objectives
- Eliminate application crashes
- Establish proper error handling patterns
- Fix context usage issues
- Create basic configuration management

### 📦 Deliverables

#### 1.1 Error Handling Overhaul
**Priority**: CRITICAL
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Create custom error types package (`internal/errors`)
  ```go
  type AppError struct {
    Code    string
    Message string
    Cause   error
  }
  ```
- [ ] Replace all `log.Fatal` and `log.Panic` with proper error returns
- [ ] Add error wrapping with context information
- [ ] Implement error recovery in critical paths

**Files to Modify**:
- `internal/adapters/google/spreadsheet/client.go` (4 fatal errors)
- `internal/adapters/telegram/handler.go` (2 panic calls)
- `internal/adapters/gemini/gemini.go` (1 fatal error)
- `cmd/telebot/main.go` (1 fatal error)

#### 1.2 Context Management Fix
**Priority**: HIGH
**Estimated Time**: 1-2 days

**Tasks**:
- [ ] Replace `context.TODO()` with proper context propagation
- [ ] Add timeout configurations for AI operations
- [ ] Implement context cancellation support

**Files to Modify**:
- `internal/adapters/telegram/handler.go:155,189`
- Add context parameter to service interfaces

#### 1.3 Basic Configuration Management
**Priority**: HIGH
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Create `internal/config` package
- [ ] Define configuration struct with validation
- [ ] Centralize environment variable loading
- [ ] Add configuration validation and defaults

**New Files**:
```
internal/config/
├── config.go
├── validation.go
└── env.go
```

**Configuration Structure**:
```go
type Config struct {
    Telegram    TelegramConfig
    Gemini      GeminiConfig
    GoogleSheets GoogleSheetsConfig
    App         AppConfig
}
```

### 🧪 Testing Requirements
- [ ] Add error handling test cases
- [ ] Test configuration loading with various scenarios
- [ ] Add context cancellation tests

### ✅ Success Criteria
- ✅ Zero `log.Fatal` or `log.Panic` calls in codebase
- ✅ All operations return proper errors
- ✅ Configuration loaded centrally with validation
- ✅ Context properly propagated through call chain

---

## 🏗️ Phase 2: Architecture & Structure (Week 3-4)
**Goal**: Establish clean architecture patterns and improve code organization.

### 🎯 Objectives
- Implement proper dependency injection
- Create validation layer
- Establish constants and eliminate magic numbers
- Improve interface design

### 📦 Deliverables

#### 2.1 Dependency Injection Container
**Priority**: HIGH
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Create dependency injection container
- [ ] Define service interfaces clearly
- [ ] Implement factory pattern for service creation
- [ ] Add interface compliance tests

**New Package Structure**:
```
internal/container/
├── container.go
├── interfaces.go
└── factory.go
```

#### 2.2 Validation Layer
**Priority**: HIGH
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Create centralized validation package
- [ ] Add input sanitization for user messages
- [ ] Implement transaction data validation
- [ ] Add amount parsing and validation utilities

**New Files**:
```
internal/validation/
├── transaction.go
├── input.go
├── amount.go
└── sanitization.go
```

#### 2.3 Constants & Configuration
**Priority**: MEDIUM
**Estimated Time**: 2 days

**Tasks**:
- [ ] Extract all magic numbers to constants
- [ ] Create business rules configuration
- [ ] Add API limits and timeout constants
- [ ] Move predefined lists to configuration

**Constants to Extract**:
- Timeout values (60 seconds)
- Sheet ranges ("detailed!A:G", "summary!A2:F12")
- Timezone strings ("Asia/Bangkok")
- File paths and directories
- API limits and retry counts

#### 2.4 Enhanced Interfaces
**Priority**: MEDIUM
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Improve existing interfaces for better testability
- [ ] Add repository pattern interfaces
- [ ] Create service layer interfaces
- [ ] Add adapter interfaces for external services

### 🧪 Testing Requirements
- [ ] Add validation test suite
- [ ] Test dependency injection scenarios
- [ ] Add interface compliance tests
- [ ] Test configuration edge cases

### ✅ Success Criteria
- ✅ Clean dependency injection throughout application
- ✅ Comprehensive input validation
- ✅ Zero magic numbers in codebase
- ✅ Well-defined interfaces for all services

---

## 🧪 Phase 3: Testing & Quality (Week 5-6)
**Goal**: Achieve comprehensive test coverage and improve code quality.

### 🎯 Objectives
- Increase test coverage to >80%
- Add integration tests
- Improve mock implementations
- Add performance benchmarks

### 📦 Deliverables

#### 3.1 Comprehensive Unit Tests
**Priority**: HIGH
**Estimated Time**: 4-5 days

**Tasks**:
- [ ] Add unit tests for all business logic
- [ ] Create table-driven tests for data validation
- [ ] Test error scenarios comprehensively
- [ ] Add edge case testing

**Test Coverage Targets**:
- `internal/service/transactions`: 95%
- `internal/domain/transactions`: 90%
- `internal/adapters/gemini`: 85%
- `internal/adapters/telegram`: 80%
- `internal/adapters/google/spreadsheet`: 80%

#### 3.2 Integration Tests
**Priority**: HIGH
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Create end-to-end workflow tests
- [ ] Add database integration tests
- [ ] Test AI integration scenarios
- [ ] Add file upload/download tests

**Integration Test Scenarios**:
- Complete image processing workflow
- Text message processing workflow
- Error recovery scenarios
- Configuration loading scenarios

#### 3.3 Enhanced Mocks & Test Utilities
**Priority**: MEDIUM
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Improve mock implementations with validation
- [ ] Create test fixtures and helpers
- [ ] Add test data generators
- [ ] Create testing utilities package

**New Testing Structure**:
```
internal/testing/
├── mocks/
├── fixtures/
├── helpers/
└── generators/
```

#### 3.4 Performance & Benchmark Tests
**Priority**: LOW
**Estimated Time**: 2 days

**Tasks**:
- [ ] Add benchmark tests for critical operations
- [ ] Create performance regression tests
- [ ] Add memory usage tests
- [ ] Create load testing scenarios

### 🧪 Testing Requirements
- [ ] All new tests must pass
- [ ] Test coverage reports generated
- [ ] Integration tests run in CI/CD
- [ ] Performance benchmarks established

### ✅ Success Criteria
- ✅ >80% test coverage across all packages
- ✅ Comprehensive integration test suite
- ✅ High-quality mock implementations
- ✅ Performance benchmarks established

---

## 💾 Phase 4: Data & Persistence (Week 7-8)
**Goal**: Improve data handling, implement proper types, and add repository pattern.

### 🎯 Objectives
- Create proper Money/Amount type
- Implement repository pattern
- Add data validation and migration utilities
- Improve file management

### 📦 Deliverables

#### 4.1 Money/Amount Type System
**Priority**: HIGH
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Create dedicated Money type with decimal precision
- [ ] Add currency conversion utilities
- [ ] Implement amount formatting standardization
- [ ] Add validation for amount ranges and formats

**New Package**:
```go
internal/types/money/
├── money.go          // Core Money type
├── currency.go       // Currency definitions
├── formatting.go     // Display formatting
└── validation.go     // Amount validation
```

#### 4.2 Repository Pattern Implementation
**Priority**: HIGH
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Create repository interfaces
- [ ] Implement Google Sheets repository
- [ ] Add transaction patterns for data consistency
- [ ] Create data access layer abstraction

**Repository Structure**:
```
internal/repository/
├── interfaces.go
├── transaction.go
├── spreadsheet/
└── memory/          // For testing
```

#### 4.3 File Management Service
**Priority**: MEDIUM
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Create dedicated file management service
- [ ] Implement proper cleanup with defer statements
- [ ] Add file validation (size, type, security)
- [ ] Create temporary file management utilities

**File Service Structure**:
```
internal/services/file/
├── service.go
├── validation.go
├── cleanup.go
└── temp.go
```

#### 4.4 Data Migration & Utilities
**Priority**: LOW
**Estimated Time**: 2 days

**Tasks**:
- [ ] Add data migration utilities for schema changes
- [ ] Create data backup and restore utilities
- [ ] Add data consistency checks
- [ ] Implement data export/import features

### 🧪 Testing Requirements
- [ ] Test Money type operations
- [ ] Test repository implementations
- [ ] Test file management operations
- [ ] Test data migration scenarios

### ✅ Success Criteria
- ✅ Robust Money type with proper decimal handling
- ✅ Repository pattern implemented for all data operations
- ✅ Secure and efficient file management
- ✅ Data consistency and migration utilities

---

## 🚀 Phase 5: Performance & Monitoring (Week 9-10)
**Goal**: Optimize performance, add monitoring, and implement production-ready features.

### 🎯 Objectives
- Add monitoring and observability
- Implement performance optimizations
- Add security enhancements
- Create production deployment features

### 📦 Deliverables

#### 5.1 Monitoring & Observability
**Priority**: HIGH
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Implement structured logging with levels
- [ ] Add metrics collection (processing time, success/failure rates)
- [ ] Create health check endpoints
- [ ] Add distributed tracing support

**Monitoring Structure**:
```
internal/monitoring/
├── logger.go
├── metrics.go
├── health.go
└── tracing.go
```

#### 5.2 Performance Optimizations
**Priority**: MEDIUM
**Estimated Time**: 3-4 days

**Tasks**:
- [ ] Add concurrent processing for independent operations
- [ ] Implement connection pooling for external services
- [ ] Add caching for frequently accessed data
- [ ] Optimize AI API call batching

**Performance Features**:
- Connection pooling for Google Sheets API
- Caching for category summaries
- Batch processing for multiple transactions
- Async file processing

#### 5.3 Security Enhancements
**Priority**: HIGH
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Add comprehensive input sanitization
- [ ] Implement file upload security checks
- [ ] Add rate limiting for API calls
- [ ] Create request validation middleware

**Security Features**:
- File type validation and scanning
- Rate limiting per user
- Input sanitization for all user data
- API key rotation support

#### 5.4 Production Features
**Priority**: MEDIUM
**Estimated Time**: 2-3 days

**Tasks**:
- [ ] Implement graceful shutdown handling
- [ ] Add configuration for different environments
- [ ] Create deployment scripts and documentation
- [ ] Add operational runbooks

**Production Features**:
- Graceful shutdown with cleanup
- Environment-specific configurations
- Docker containerization improvements
- Kubernetes deployment manifests

### 🧪 Testing Requirements
- [ ] Load testing for performance optimizations
- [ ] Security testing for input validation
- [ ] Monitoring system testing
- [ ] Production deployment testing

### ✅ Success Criteria
- ✅ Comprehensive monitoring and alerting
- ✅ Significant performance improvements
- ✅ Enhanced security posture
- ✅ Production-ready deployment

---

## 🗓️ Timeline Summary

| Phase | Duration | Key Focus | Team Size |
|-------|----------|-----------|-----------|
| Phase 1 | Week 1-2 | Stability & Foundation | 2-3 developers |
| Phase 2 | Week 3-4 | Architecture & Structure | 2-3 developers |
| Phase 3 | Week 5-6 | Testing & Quality | 1-2 developers |
| Phase 4 | Week 7-8 | Data & Persistence | 2 developers |
| Phase 5 | Week 9-10 | Performance & Monitoring | 1-2 developers |

**Total Duration**: 10 weeks
**Recommended Team**: 2-3 developers with Go experience

## 🔄 Phase Dependencies

```
Phase 1 (Foundation)
    ↓
Phase 2 (Architecture) ← Must complete Phase 1 first
    ↓
Phase 3 (Testing) ← Can start after Phase 2
    ↓
Phase 4 (Data) ← Can partially overlap with Phase 3
    ↓
Phase 5 (Performance) ← Can start after Phase 4 foundation
```

## 📊 Risk Mitigation

### High Risk Areas
1. **Phase 1 Error Handling**: Risk of breaking existing functionality
   - **Mitigation**: Comprehensive testing at each step
2. **Phase 2 Architecture**: Risk of over-engineering
   - **Mitigation**: Focus on immediate needs, avoid premature optimization
3. **Phase 4 Data Types**: Risk of data corruption during migration
   - **Mitigation**: Thorough backup and rollback procedures

### Rollback Strategy
Each phase should include:
- Feature flags for new functionality
- Ability to rollback to previous phase
- Comprehensive testing before promotion
- Gradual rollout with monitoring

---

*This plan preserves all existing functionality while significantly improving code quality, maintainability, and reliability. Each improvement can be implemented incrementally without breaking changes. The phased approach ensures steady progress while maintaining system stability and allows for course correction at each phase boundary.*