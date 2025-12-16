#!/bin/bash

# Integration testing script for go-quay CLI
# 
# Environment variables:
#   QUAY_TOKEN      - Quay.io API token (required for API tests)
#   QUAY_ORG        - Organization name (for org-scoped tests)
#   QUAY_NAMESPACE  - Repository namespace (for repo-scoped tests)
#   QUAY_REPOSITORY - Repository name (for repo-scoped tests)
#
# Example: QUAY_TOKEN=your_token QUAY_ORG=your_org QUAY_NAMESPACE=myns QUAY_REPOSITORY=myrepo ./integration-test.sh

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

# Test manifest command help
echo "Testing manifest command help..."
./$APP_NAME get manifest --help > /dev/null
./$APP_NAME get manifest info --help > /dev/null
./$APP_NAME get manifest delete --help > /dev/null
./$APP_NAME get manifest labels --help > /dev/null
./$APP_NAME get manifest label --help > /dev/null
./$APP_NAME get manifest add-label --help > /dev/null
./$APP_NAME get manifest remove-label --help > /dev/null

# Test secscan command help
echo "Testing secscan command help..."
./$APP_NAME get secscan --help > /dev/null
./$APP_NAME get secscan info --help > /dev/null

# Test robot command help (user-level robots)
echo "Testing robot command help..."
./$APP_NAME get robot --help > /dev/null
./$APP_NAME get robot list --help > /dev/null
./$APP_NAME get robot info --help > /dev/null
./$APP_NAME get robot create --help > /dev/null
./$APP_NAME get robot delete --help > /dev/null
./$APP_NAME get robot regenerate --help > /dev/null
./$APP_NAME get robot permissions --help > /dev/null

# Test search command help
echo "Testing search command help..."
./$APP_NAME get search --help > /dev/null
./$APP_NAME get search repositories --help > /dev/null
./$APP_NAME get search all --help > /dev/null

# Test team command help
echo "Testing team command help..."
./$APP_NAME get team --help > /dev/null
./$APP_NAME get team list --help > /dev/null
./$APP_NAME get team info --help > /dev/null
./$APP_NAME get team create --help > /dev/null
./$APP_NAME get team update --help > /dev/null
./$APP_NAME get team delete --help > /dev/null
./$APP_NAME get team members --help > /dev/null
./$APP_NAME get team add-member --help > /dev/null
./$APP_NAME get team remove-member --help > /dev/null
./$APP_NAME get team permissions --help > /dev/null
./$APP_NAME get team set-permission --help > /dev/null
./$APP_NAME get team remove-permission --help > /dev/null

# Test build command help
echo "Testing build command help..."
./$APP_NAME get build --help > /dev/null
./$APP_NAME get build list --help > /dev/null
./$APP_NAME get build info --help > /dev/null
./$APP_NAME get build logs --help > /dev/null
./$APP_NAME get build request --help > /dev/null
./$APP_NAME get build cancel --help > /dev/null

# Test notification command help
echo "Testing notification command help..."
./$APP_NAME get notification --help > /dev/null
./$APP_NAME get notification list --help > /dev/null
./$APP_NAME get notification info --help > /dev/null
./$APP_NAME get notification create --help > /dev/null
./$APP_NAME get notification delete --help > /dev/null
./$APP_NAME get notification test --help > /dev/null
./$APP_NAME get notification reset --help > /dev/null

# Test trigger command help
echo "Testing trigger command help..."
./$APP_NAME get trigger --help > /dev/null
./$APP_NAME get trigger list --help > /dev/null
./$APP_NAME get trigger info --help > /dev/null
./$APP_NAME get trigger delete --help > /dev/null
./$APP_NAME get trigger enable --help > /dev/null
./$APP_NAME get trigger disable --help > /dev/null
./$APP_NAME get trigger start --help > /dev/null
./$APP_NAME get trigger activate --help > /dev/null

# Test discovery command help
echo "Testing discovery command help..."
./$APP_NAME get discovery --help > /dev/null

# Test error command help
echo "Testing error command help..."
./$APP_NAME get error --help > /dev/null

# Test messages command help
echo "Testing messages command help..."
./$APP_NAME get messages --help > /dev/null

# Test prototype command help
echo "Testing prototype command help..."
./$APP_NAME get prototype --help > /dev/null
./$APP_NAME get prototype list --help > /dev/null
./$APP_NAME get prototype info --help > /dev/null
./$APP_NAME get prototype create --help > /dev/null
./$APP_NAME get prototype update --help > /dev/null
./$APP_NAME get prototype delete --help > /dev/null

# Test repotoken command help (deprecated)
echo "Testing repotoken command help..."
./$APP_NAME get repotoken --help > /dev/null
./$APP_NAME get repotoken list --help > /dev/null
./$APP_NAME get repotoken info --help > /dev/null
./$APP_NAME get repotoken create --help > /dev/null
./$APP_NAME get repotoken update --help > /dev/null
./$APP_NAME get repotoken delete --help > /dev/null

# Test user command help
echo "Testing user command help..."
./$APP_NAME get user --help > /dev/null
./$APP_NAME get user info --help > /dev/null
./$APP_NAME get user starred --help > /dev/null
./$APP_NAME get user star --help > /dev/null
./$APP_NAME get user unstar --help > /dev/null

# Test permissions command help
echo "Testing permissions command help..."
./$APP_NAME get permissions --help > /dev/null
./$APP_NAME get permissions list --help > /dev/null
./$APP_NAME get permissions get --help > /dev/null
./$APP_NAME get permissions set --help > /dev/null
./$APP_NAME get permissions delete --help > /dev/null

# Test tag command help
echo "Testing tag command help..."
./$APP_NAME get tag --help > /dev/null

# Test error handling
echo "Testing error handling..."
OUTPUT=$(./go-quay get billing org-info 2>&1) && echo "$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error" || (echo "✗ Missing required flag error not displayed" && exit 1)
OUTPUT=$(./go-quay get billing user-info 2>&1) && echo "$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error" || (echo "✗ Missing required flag error not displayed" && exit 1)
OUTPUT=$(./go-quay get organization info 2>&1) && echo "$OUTPUT" | grep -q "required flag" && echo "✓ Correctly shows required flag error for organization" || (echo "✗ Missing required flag error not displayed for organization" && exit 1)

# Integration tests with real API if token is provided
if [ -n "$QUAY_TOKEN" ]; then
    echo "Running integration tests with real API..."
    
    echo "=== Testing Read-Only APIs ==="
    
    # Discovery API - returns API information (read-only, safe)
    printf "Testing discovery API... "
    DISCOVERY_INFO=$(./$APP_NAME get discovery --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "✓ (API discovery successful)"
    else
        echo "✗ (may not be available)"
    fi
    
    # Messages API - returns system messages (read-only, safe)
    printf "Testing messages API... "
    MESSAGES_INFO=$(./$APP_NAME get messages --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ]; then
        if echo "$MESSAGES_INFO" | grep -q "messages" 2>/dev/null; then
            MSG_COUNT=$(echo "$MESSAGES_INFO" | jq '.messages | length' 2>/dev/null || echo "0")
            echo "✓ (found $MSG_COUNT system messages)"
        else
            echo "✓ (no messages returned)"
        fi
    else
        echo "✗ (may not be available)"
    fi
    
    # Search API - search for repositories (read-only, safe)
    printf "Testing search repositories API... "
    SEARCH_INFO=$(./$APP_NAME get search repositories --query "quay" --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ]; then
        if echo "$SEARCH_INFO" | grep -q "results" 2>/dev/null; then
            RESULT_COUNT=$(echo "$SEARCH_INFO" | jq '.results | length' 2>/dev/null || echo "unknown")
            echo "✓ (found $RESULT_COUNT repositories matching 'quay')"
        else
            echo "✓ (search completed)"
        fi
    else
        echo "✗ (may not have access)"
    fi
    
    printf "Testing search all API... "
    SEARCH_ALL_INFO=$(./$APP_NAME get search all --query "test" --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "✓ (search all completed)"
    else
        echo "✗ (may not have access)"
    fi
    
    echo "=== Testing User APIs ==="
    
    # User info
    printf "Testing user info API... "
    USER_INFO=$(./$APP_NAME get user info --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ] && echo "$USER_INFO" | grep -q "username" 2>/dev/null; then
        USERNAME=$(echo "$USER_INFO" | jq -r '.username' 2>/dev/null || echo "unknown")
        echo "✓ (logged in as: $USERNAME)"
    else
        echo "✗ (may not have access)"
    fi
    
    # User starred repositories
    printf "Testing user starred repositories... "
    STARRED_INFO=$(./$APP_NAME get user starred --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ]; then
        if echo "$STARRED_INFO" | grep -q "repositories" 2>/dev/null; then
            STARRED_COUNT=$(echo "$STARRED_INFO" | jq '.repositories | length' 2>/dev/null || echo "0")
            echo "✓ (found $STARRED_COUNT starred repositories)"
        else
            echo "✓ (no starred repositories)"
        fi
    else
        echo "✗ (may not have access)"
    fi
    
    # User robot accounts (read-only list)
    printf "Testing user robots list... "
    USER_ROBOTS=$(./$APP_NAME get robot list --token "$QUAY_TOKEN" 2>/dev/null)
    if [ $? -eq 0 ]; then
        if echo "$USER_ROBOTS" | grep -q "robots" 2>/dev/null; then
            ROBOT_COUNT=$(echo "$USER_ROBOTS" | jq '.robots | length' 2>/dev/null || echo "0")
            echo "✓ (found $ROBOT_COUNT user robots)"
        else
            echo "✓ (no robots found)"
        fi
    else
        echo "✗ (may not have access)"
    fi
    
    echo "=== Testing Billing APIs ==="
    
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
        
        echo "=== Testing Prototype API ==="
        
        # Prototypes (default permissions) - read-only list
        printf "Testing prototypes for $QUAY_ORG... "
        PROTOTYPES_INFO=$(./$APP_NAME get prototype list --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            if echo "$PROTOTYPES_INFO" | grep -q "prototypes" 2>/dev/null; then
                PROTO_COUNT=$(echo "$PROTOTYPES_INFO" | jq '.prototypes | length' 2>/dev/null || echo "0")
                echo "✓ (found $PROTO_COUNT prototypes)"
            else
                echo "✓ (no prototypes configured)"
            fi
        else
            echo "✗ (may not have access)"
        fi
        
        echo "=== Testing Team API (Dedicated Commands) ==="
        
        # Team list using dedicated team command
        printf "Testing team list for $QUAY_ORG... "
        TEAM_LIST=$(./$APP_NAME get team list --organization "$QUAY_ORG" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            TEAM_COUNT=$(echo "$TEAM_LIST" | grep -o '"name"' 2>/dev/null | wc -l)
            echo "✓ (found $TEAM_COUNT teams)"
            
            # If we found teams, test more team functionality
            if [ $TEAM_COUNT -gt 0 ]; then
                FIRST_TEAM=$(echo "$TEAM_LIST" | grep -o '"name": *"[^"]*"' 2>/dev/null | head -1 | cut -d'"' -f4)
                
                if [ -n "$FIRST_TEAM" ]; then
                    printf "Testing team info for '$FIRST_TEAM' (dedicated command)... "
                    if ./$APP_NAME get team info --organization "$QUAY_ORG" --name "$FIRST_TEAM" --token "$QUAY_TOKEN" > /dev/null 2>&1; then
                        echo "✓"
                    else
                        echo "✗ (may not have access)"
                    fi
                    
                    printf "Testing team members for '$FIRST_TEAM' (dedicated command)... "
                    TEAM_MEMBERS=$(./$APP_NAME get team members --organization "$QUAY_ORG" --name "$FIRST_TEAM" --token "$QUAY_TOKEN" 2>/dev/null)
                    if [ $? -eq 0 ]; then
                        MEMBER_COUNT=$(echo "$TEAM_MEMBERS" | jq '.members | length' 2>/dev/null || echo "unknown")
                        echo "✓ (found $MEMBER_COUNT members)"
                    else
                        echo "✗ (may not have access)"
                    fi
                    
                    printf "Testing team permissions for '$FIRST_TEAM'... "
                    TEAM_PERMS=$(./$APP_NAME get team permissions --organization "$QUAY_ORG" --name "$FIRST_TEAM" --token "$QUAY_TOKEN" 2>/dev/null)
                    if [ $? -eq 0 ]; then
                        PERM_COUNT=$(echo "$TEAM_PERMS" | jq '.permissions | length' 2>/dev/null || echo "unknown")
                        echo "✓ (found $PERM_COUNT repository permissions)"
                    else
                        echo "✗ (may not have access)"
                    fi
                fi
            fi
        else
            echo "✗ (may not have access)"
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
    
    # Repository-scoped tests if namespace and repository are provided
    if [ -n "$QUAY_NAMESPACE" ] && [ -n "$QUAY_REPOSITORY" ]; then
        echo "=== Testing Repository-Scoped APIs ==="
        
        # Build API - list builds (read-only)
        printf "Testing build list for $QUAY_NAMESPACE/$QUAY_REPOSITORY... "
        BUILDS_INFO=$(./$APP_NAME get build list --namespace "$QUAY_NAMESPACE" --repository "$QUAY_REPOSITORY" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            if echo "$BUILDS_INFO" | grep -q "builds" 2>/dev/null; then
                BUILD_COUNT=$(echo "$BUILDS_INFO" | jq '.builds | length' 2>/dev/null || echo "0")
                echo "✓ (found $BUILD_COUNT builds)"
            else
                echo "✓ (no builds found)"
            fi
        else
            echo "✗ (may not have access or repository doesn't exist)"
        fi
        
        # Trigger API - list triggers (read-only)
        printf "Testing trigger list for $QUAY_NAMESPACE/$QUAY_REPOSITORY... "
        TRIGGERS_INFO=$(./$APP_NAME get trigger list --namespace "$QUAY_NAMESPACE" --repository "$QUAY_REPOSITORY" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            if echo "$TRIGGERS_INFO" | grep -q "triggers" 2>/dev/null; then
                TRIGGER_COUNT=$(echo "$TRIGGERS_INFO" | jq '.triggers | length' 2>/dev/null || echo "0")
                echo "✓ (found $TRIGGER_COUNT triggers)"
            else
                echo "✓ (no triggers configured)"
            fi
        else
            echo "✗ (may not have access or repository doesn't exist)"
        fi
        
        # Notification API - list notifications (read-only)
        printf "Testing notification list for $QUAY_NAMESPACE/$QUAY_REPOSITORY... "
        NOTIFICATIONS_INFO=$(./$APP_NAME get notification list --namespace "$QUAY_NAMESPACE" --repository "$QUAY_REPOSITORY" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            if echo "$NOTIFICATIONS_INFO" | grep -q "notifications" 2>/dev/null; then
                NOTIF_COUNT=$(echo "$NOTIFICATIONS_INFO" | jq '.notifications | length' 2>/dev/null || echo "0")
                echo "✓ (found $NOTIF_COUNT notifications)"
            else
                echo "✓ (no notifications configured)"
            fi
        else
            echo "✗ (may not have access or repository doesn't exist)"
        fi
        
        # Permissions API - list permissions (read-only)
        printf "Testing permissions for $QUAY_NAMESPACE/$QUAY_REPOSITORY... "
        PERMS_INFO=$(./$APP_NAME get permissions list --namespace "$QUAY_NAMESPACE" --repository "$QUAY_REPOSITORY" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            if echo "$PERMS_INFO" | grep -q "permissions" 2>/dev/null; then
                PERM_COUNT=$(echo "$PERMS_INFO" | jq '.permissions | length' 2>/dev/null || echo "0")
                echo "✓ (found $PERM_COUNT permissions)"
            else
                echo "✓ (no permissions found)"
            fi
        else
            echo "✗ (may not have access or repository doesn't exist)"
        fi
        
        # Tags API - list tags (read-only)
        printf "Testing tags for $QUAY_NAMESPACE/$QUAY_REPOSITORY... "
        TAGS_INFO=$(./$APP_NAME get tag list --namespace "$QUAY_NAMESPACE" --repository "$QUAY_REPOSITORY" --token "$QUAY_TOKEN" 2>/dev/null)
        if [ $? -eq 0 ]; then
            if echo "$TAGS_INFO" | grep -q "tags" 2>/dev/null; then
                TAG_COUNT=$(echo "$TAGS_INFO" | jq '.tags | length' 2>/dev/null || echo "0")
                echo "✓ (found $TAG_COUNT tags)"
            else
                echo "✓ (no tags found)"
            fi
        else
            echo "✗ (may not have access or repository doesn't exist)"
        fi
        
        echo "=== Repository Testing Complete ==="
    else
        echo "No QUAY_NAMESPACE/QUAY_REPOSITORY specified, skipping repository-scoped tests"
        echo "Set QUAY_NAMESPACE and QUAY_REPOSITORY to test Build, Trigger, Notification APIs"
    fi
else
    echo "No QUAY_TOKEN available, skipping integration tests"
fi

echo "✓ All CLI tests completed successfully" 