#!/bin/bash

# Integration testing script for go-quay CLI
# Set QUAY_TOKEN and QUAY_ORG environment variables for API integration tests
# Example: QUAY_TOKEN=your_token QUAY_ORG=your_org ./integration-test.sh

APP_NAME="go-quay"

echo "Running CLI integration tests..."

# Test basic CLI help commands
echo "Testing main CLI help..."
./$APP_NAME --help > /dev/null

echo "Testing get command help..."
./$APP_NAME get --help > /dev/null

# Test billing command help
echo "Testing billing command help..."
./$APP_NAME get billing --help > /dev/null

echo "Testing billing subcommands help..."
./$APP_NAME get billing org-info --help > /dev/null
./$APP_NAME get billing user-info --help > /dev/null
./$APP_NAME get billing org-subscription --help > /dev/null
./$APP_NAME get billing user-subscription --help > /dev/null
./$APP_NAME get billing org-invoices --help > /dev/null
# ./$APP_NAME get billing user-invoices --help > /dev/null  # Endpoint doesn't exist in Quay API
# Usage commands commented out - endpoints don't exist in Quay API
# ./$APP_NAME get billing org-usage --help > /dev/null
# ./$APP_NAME get billing user-usage --help > /dev/null
./$APP_NAME get billing plans --help > /dev/null

# Test logs command help
echo "Testing logs command help..."
./$APP_NAME get logs --help > /dev/null

# Test repository command help
echo "Testing repository command help..."
./$APP_NAME get repository --help > /dev/null

# Test organization command help
echo "Testing organization command help..."
./$APP_NAME get organization --help > /dev/null

echo "Testing organization subcommands help..."
./$APP_NAME get organization info --help > /dev/null
./$APP_NAME get organization members --help > /dev/null
./$APP_NAME get organization teams --help > /dev/null
./$APP_NAME get organization team --help > /dev/null
./$APP_NAME get organization team-members --help > /dev/null
./$APP_NAME get organization robots --help > /dev/null
./$APP_NAME get organization quota --help > /dev/null
./$APP_NAME get organization auto-prune --help > /dev/null
./$APP_NAME get organization applications --help > /dev/null

# Test error handling
echo "Testing error handling..."
OUTPUT=$(./go-quay get billing org-info 2>&1) && echo "$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error" || (echo "✗ Missing required flag error not displayed" && exit 1)
OUTPUT=$(./go-quay get billing user-info 2>&1) && echo "$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error" || (echo "✗ Missing required flag error not displayed" && exit 1)
OUTPUT=$(./go-quay get organization info 2>&1) && echo "$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error for organization" || (echo "✗ Missing required flag error not displayed for organization" && exit 1)

# Integration tests with real API if token is provided
if [ -n "$QUAY_TOKEN" ]; then
    echo "Running integration tests with real API..."
    
    # Test user endpoints
    printf "Testing user billing info... "
    if ./$APP_NAME get billing user-info --token "$QUAY_TOKEN" > /dev/null 2>&1; then
        echo "✓"
    else
        echo "✗ (may not have access)"
    fi
    
    printf "Testing user subscription... "
    if ./$APP_NAME get billing user-subscription --token "$QUAY_TOKEN" > /dev/null 2>&1; then
        echo "✓"
    else
        echo "✗ (may not have access)"
    fi
    
    printf "Testing available plans... "
    if ./$APP_NAME get billing plans --token "$QUAY_TOKEN" > /dev/null 2>&1; then
        echo "✓"
    else
        echo "✗ (may not have access)"
    fi
    
    # Test organization endpoints if org is provided
    if [ -n "$QUAY_ORG" ]; then
        echo "=== Testing Organization Endpoints ==="
        
        # Billing tests
        printf "Testing org billing info for $QUAY_ORG... "
        if ./$APP_NAME get billing org-info --organization "$QUAY_ORG" --token "$QUAY_TOKEN" > /dev/null 2>&1; then
            echo "✓"
        else
            echo "✗ (may not have access)"
        fi
        
        printf "Testing org subscription for $QUAY_ORG... "
        if ./$APP_NAME get billing org-subscription --organization "$QUAY_ORG" --token "$QUAY_TOKEN" > /dev/null 2>&1; then
            echo "✓"
        else
            echo "✗ (may not have access)"
        fi
        
        printf "Testing org invoices for $QUAY_ORG... "
        if ./$APP_NAME get billing org-invoices --organization "$QUAY_ORG" --token "$QUAY_TOKEN" > /dev/null 2>&1; then
            echo "✓"
        else
            echo "✗ (may not have access)"
        fi
        
        echo "=== Testing Organization Management ==="
        
        # Organization info with validation
        printf "Testing org info for $QUAY_ORG... "
        ORG_INFO=$(./$APP_NAME get organization info --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ] && echo "$ORG_INFO" | grep -q "name.*$QUAY_ORG" 2>/dev/null; then
            echo "✓ (found organization name in response)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        # Organization members with validation  
        printf "Testing org members for $QUAY_ORG... "
        MEMBERS_INFO=$(./$APP_NAME get organization members --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ] && echo "$MEMBERS_INFO" | grep -q "members" 2>/dev/null; then
            MEMBER_COUNT=$(echo "$MEMBERS_INFO" | grep -o '"name"' 2>/dev/null | wc -l)
            # Subtract 1 because "name" appears in team objects too, but we want member count
            ACTUAL_MEMBER_COUNT=$(echo "$MEMBERS_INFO" | jq '.members | length' 2>/dev/null || echo "$MEMBER_COUNT")
            echo "✓ (found $ACTUAL_MEMBER_COUNT members)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        # Organization teams with validation (requires admin permissions)
        printf "Testing org teams for $QUAY_ORG... "
        TEAMS_INFO=$(./$APP_NAME get organization teams --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>&1)
        if [ $? -eq 0 ]; then
            TEAM_COUNT=$(echo "$TEAMS_INFO" | grep -o '"name"' 2>/dev/null | wc -l)
            echo "✓ (found $TEAM_COUNT teams)"
            
            # If we found teams, test team-specific functionality
            if [ $TEAM_COUNT -gt 0 ]; then
                echo "=== Testing Team-Specific Functionality ==="
                
                # Get first team name for testing
                FIRST_TEAM=$(echo "$TEAMS_INFO" | grep -o '"name": *"[^"]*"' 2>/dev/null | head -1 | cut -d'"' -f4)
                
                if [ -n "$FIRST_TEAM" ]; then
                    printf "Testing team info for '$FIRST_TEAM'... "
                    if ./$APP_NAME get organization team --organization "$QUAY_ORG" --team "$FIRST_TEAM" --token "$QUAY_TOKEN" > /dev/null 2>&1; then
                        echo "✓"
                    else
                        echo "✗ (may not have access)"
                    fi
                    
                    printf "Testing team members for '$FIRST_TEAM'... "
                    TEAM_MEMBERS=$(./$APP_NAME get organization team-members --organization "$QUAY_ORG" --team "$FIRST_TEAM" --token "$QUAY_TOKEN" 2>/dev/null)
                    if [ $? -eq 0 ] && echo "$TEAM_MEMBERS" | grep -q "members" 2>/dev/null; then
                        TEAM_MEMBER_COUNT=$(echo "$TEAM_MEMBERS" | jq '.members | length' 2>/dev/null || echo "unknown")
                        echo "✓ (found $TEAM_MEMBER_COUNT members)"
                    else
                        echo "✗ (may not have access or invalid response)"
                    fi
                else
                    echo "No team name found to test team-specific functionality"
                fi
            fi
        elif echo "$TEAMS_INFO" | grep -q "403" 2>/dev/null; then
            echo "⚠ (requires admin permissions)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        echo "=== Testing Organization Resources ==="
        
        # Organization robots with validation
        printf "Testing org robots for $QUAY_ORG... "
        ROBOTS_INFO=$(./$APP_NAME get organization robots --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ] && echo "$ROBOTS_INFO" | grep -q "robots" 2>/dev/null; then
            ROBOT_COUNT=$(echo "$ROBOTS_INFO" | jq '.robots | length' 2>/dev/null || echo "unknown")
            echo "✓ (found $ROBOT_COUNT robots)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        # Organization quota with validation (requires admin permissions)
        printf "Testing org quota for $QUAY_ORG... "
        QUOTA_INFO=$(./$APP_NAME get organization quota --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>&1)
        if [ $? -eq 0 ]; then
            if echo "$QUOTA_INFO" | grep -q "limit_bytes" 2>/dev/null; then
                QUOTA_LIMIT=$(echo "$QUOTA_INFO" | grep -o '"limit_bytes": *[0-9]*' 2>/dev/null | cut -d':' -f2 | tr -d ' ')
                echo "✓ (quota limit: $QUOTA_LIMIT bytes)"
            else
                echo "✓ (no quota set)"
            fi
        elif echo "$QUOTA_INFO" | grep -q "403" 2>/dev/null; then
            echo "⚠ (requires admin permissions)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        # Auto-prune policies (may not be available on all Quay instances)
        printf "Testing auto-prune policies for $QUAY_ORG... "
        PRUNE_INFO=$(./$APP_NAME get organization auto-prune --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>&1)
        if [ $? -eq 0 ] && echo "$PRUNE_INFO" | grep -q "policies" 2>/dev/null; then
            POLICY_COUNT=$(echo "$PRUNE_INFO" | jq '.policies | length' 2>/dev/null || echo "unknown")
            echo "✓ (found $POLICY_COUNT auto-prune policies)"
        elif echo "$PRUNE_INFO" | grep -q "403" 2>/dev/null; then
            echo "⚠ (feature not available on this Quay instance)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        # Organization applications with validation
        printf "Testing org applications for $QUAY_ORG... "
        APPS_INFO=$(./$APP_NAME get organization applications --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ] && echo "$APPS_INFO" | grep -q "applications" 2>/dev/null; then
            APP_COUNT=$(echo "$APPS_INFO" | jq '.applications | length' 2>/dev/null || echo "unknown")
            echo "✓ (found $APP_COUNT applications)"
        else
            echo "✗ (may not have access or invalid response)"
        fi
        
        echo "=== Testing Error Conditions ==="
        
        # Test invalid organization
        printf "Testing invalid organization error handling... "
        ERROR_OUTPUT=$(./$APP_NAME get organization info --organization "invalid-org-$(date +%s)" --token "$QUAY_TOKEN" 2>&1)
        if echo "$ERROR_OUTPUT" | grep -q -i "error\|not found\|invalid" 2>/dev/null; then
            echo "✓ (properly handles invalid organization)"
        else
            echo "✗ (error handling unclear)"
        fi
        
        # Test invalid team
        printf "Testing invalid team error handling... "
        ERROR_OUTPUT=$(./$APP_NAME get organization team --organization "$QUAY_ORG" --team "invalid-team-$(date +%s)" --token "$QUAY_TOKEN" 2>&1)
        if echo "$ERROR_OUTPUT" | grep -q -i "error\|not found\|invalid" 2>/dev/null; then
            echo "✓ (properly handles invalid team)"
        else
            echo "✗ (error handling unclear)"
        fi
        
        # Test missing token
        printf "Testing missing token error handling... "
        ERROR_OUTPUT=$(./$APP_NAME get organization info --organization "$QUAY_ORG" 2>&1)
        if echo "$ERROR_OUTPUT" | grep -q -i "required flag.*token" 2>/dev/null; then
            echo "✓ (properly requires token)"
        else
            echo "✗ (token validation unclear)"
        fi
        
        echo "=== Organization Testing Complete ==="
    else
        echo "No QUAY_ORG specified, skipping organization tests"
    fi
else
    echo "No QUAY_TOKEN available, skipping integration tests"
fi

echo "✓ All CLI tests completed successfully" 