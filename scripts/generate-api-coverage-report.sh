#!/bin/bash

# Generate API coverage report for go-quay project
# This script creates a comprehensive markdown report showing API coverage status

set -e

REPORT_FILE="api-coverage-report.md"

echo "Generating API coverage report..."

# Generate main report header
cat > "$REPORT_FILE" << 'EOF'
# Quay API Coverage Report
EOF

echo "Generated on: $(date)" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# Phase 1 Implementation Status
cat >> "$REPORT_FILE" << 'EOF'
## Phase 1 Implementation Status ✅

### ✅ Repository Management (CRUD)
- [x] POST /api/v1/repository - CreateRepository()
- [x] GET /api/v1/repository/{namespace}/{repository} - GetRepository()
- [x] PUT /api/v1/repository/{namespace}/{repository} - UpdateRepository()
- [x] DELETE /api/v1/repository/{namespace}/{repository} - DeleteRepository()

### ✅ Repository Permissions
- [x] GET /api/v1/repository/{namespace}/{repository}/permissions - GetRepositoryPermissions()
- [x] PUT /api/v1/repository/{namespace}/{repository}/permissions/{username} - SetRepositoryPermission()
- [x] DELETE /api/v1/repository/{namespace}/{repository}/permissions/{username} - RemoveRepositoryPermission()

### ✅ Enhanced Tag Operations
- [x] GET /api/v1/repository/{namespace}/{repository}/tag/{tag} - GetTag()
- [x] PUT /api/v1/repository/{namespace}/{repository}/tag/{tag} - UpdateTag()
- [x] DELETE /api/v1/repository/{namespace}/{repository}/tag/{tag} - DeleteTag()
- [x] GET /api/v1/repository/{namespace}/{repository}/tag/{tag}/history - GetTagHistory()
- [x] POST /api/v1/repository/{namespace}/{repository}/tag/{tag}/revert - RevertTag()

### ✅ User Account Operations
- [x] GET /api/v1/user - GetUser()
- [x] GET /api/v1/user/starred - GetStarredRepositories()
- [x] PUT /api/v1/repository/{namespace}/{repository}/star - StarRepository()
- [x] DELETE /api/v1/repository/{namespace}/{repository}/star - UnstarRepository()

## Existing Implementation Status

### ✅ Billing & Subscriptions (Complete)
### ✅ Organization Management (Complete)
### ✅ Logging (Complete)

## Future Phases

### 🔄 Phase 2: Security & Scanning
- [ ] Vulnerability reports
- [ ] Security scan results
- [ ] Image signing

### 🔄 Phase 3: Build System
- [ ] Build triggers
- [ ] Build history and logs
- [ ] Build configuration

### 🔄 Phase 4: Advanced Features
- [ ] Repository search
- [ ] Webhooks & notifications
- [ ] Repository mirroring

EOF

# Generate test coverage statistics
echo "" >> "$REPORT_FILE"
echo "## Test Coverage Statistics" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### Unit Tests" >> "$REPORT_FILE"

# Run tests and capture coverage
echo "Running tests and generating coverage report..."
go test -v -coverprofile=coverage.out ./... > test_output.txt 2>&1 || true

if [ -f coverage.out ]; then
    echo "Processing coverage data..."
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "- Total coverage: $COVERAGE" >> "$REPORT_FILE"
    
    # Count tests by category
    echo "Counting tests by category..."
    REPO_TESTS=$(grep -c "Test.*Repository" test_output.txt || echo "0")
    PERM_TESTS=$(grep -c "Test.*Permission" test_output.txt || echo "0") 
    TAG_TESTS=$(grep -c "Test.*Tag" test_output.txt || echo "0")
    USER_TESTS=$(grep -c "Test.*User" test_output.txt || echo "0")
    
    echo "- Repository tests: $REPO_TESTS" >> "$REPORT_FILE"
    echo "- Permission tests: $PERM_TESTS" >> "$REPORT_FILE"
    echo "- Tag tests: $TAG_TESTS" >> "$REPORT_FILE"
    echo "- User tests: $USER_TESTS" >> "$REPORT_FILE"
else
    echo "- Coverage data not available" >> "$REPORT_FILE"
fi

# Clean up temporary files
echo "Cleaning up temporary files..."
rm -f test_output.txt

echo "API coverage report generated successfully: $REPORT_FILE"
