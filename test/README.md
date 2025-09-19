# Ecommerce Shop - Unit Tests

This directory contains comprehensive unit tests for the ecommerce shop application.

## Test Structure

```
test/
├── README.md                 # This file
└── testdata/                 # Test data files (if needed)

internal/
├── testutils/               # Test utilities and helpers
│   ├── mock.go             # Database mocking utilities
│   ├── helpers.go          # Test helper functions
│   └── config.go           # Test configuration
├── handlers/               # Handler tests
│   ├── auth_test.go        # Authentication handler tests
│   ├── products_test.go    # Products handler tests
│   ├── orders_test.go      # Orders handler tests
│   └── warehouses_test.go  # Warehouses handler tests
├── service/                # Service layer tests
│   ├── auth_test.go        # Authentication service tests
│   ├── products_test.go    # Products service tests
│   ├── orders_test.go      # Orders service tests
│   └── warehouses_test.go  # Warehouses service tests
├── entity/                 # Entity/Model tests
│   ├── auth_test.go        # Authentication entity tests
│   ├── products_test.go    # Products entity tests
│   ├── orders_test.go      # Orders entity tests
│   └── warehouses_test.go  # Warehouses entity tests
└── helpers/                # Helper function tests
    └── response_test.go    # Response helper tests
```

## Running Tests

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Run Only Unit Tests
```bash
make test-unit
```

### Run Specific Test Package
```bash
go test -v ./internal/handlers
go test -v ./internal/service
go test -v ./internal/entity
```

### Run Specific Test
```bash
go test -v -run TestAuthHandler_Register ./internal/handlers
```

## Test Coverage

The test suite provides comprehensive coverage for:

- **Handlers**: HTTP endpoint testing with mocked dependencies
- **Services**: Business logic testing with mocked database
- **Entities**: Data validation and structure testing
- **Helpers**: Utility function testing

## Test Utilities

### MockDB
Creates a mock database for testing without requiring a real database connection.

### TestGinContext
Creates a Gin context for testing HTTP handlers.

### TestValidator
Creates a validator instance for testing data validation.

### TestConfig
Creates test configuration with safe default values.

## Test Data

All test data is generated programmatically to avoid external dependencies. The tests use:
- Mock databases with sqlmock
- Generated test data
- Isolated test contexts

## Best Practices

1. **Isolation**: Each test is independent and doesn't affect others
2. **Mocking**: External dependencies are mocked to ensure fast, reliable tests
3. **Coverage**: Tests cover both success and error scenarios
4. **Validation**: Input validation is thoroughly tested
5. **Edge Cases**: Boundary conditions and edge cases are tested

## Dependencies

The test suite uses the following testing libraries:
- `github.com/stretchr/testify` - Assertions and test utilities
- `github.com/DATA-DOG/go-sqlmock` - Database mocking
- `go.uber.org/zap/zaptest` - Test logging

## Adding New Tests

When adding new functionality:

1. Create corresponding test files in the appropriate package
2. Use the existing test utilities in `internal/testutils`
3. Follow the naming convention: `*_test.go`
4. Include both positive and negative test cases
5. Mock external dependencies
6. Ensure tests are isolated and repeatable
