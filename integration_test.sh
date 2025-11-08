#!/bin/bash

# Integration tests for ecs-tag-shift

BINARY="./ecs-tag-shift"
EXAMPLES_DIR="./examples"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

passed=0
failed=0

echo "=== Building binary ==="
go build -o "$BINARY" ./cmd/ecs-tag-shift || { echo "Build failed"; exit 1; }
echo "Build successful"

echo ""
echo "=== Running show command tests ==="

# Test show with task definition
echo -n "Test: Show task definition (default) ... "
if $BINARY show $EXAMPLES_DIR/task-definition.json 2>&1 | grep -q "my-app"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test show with text output
echo -n "Test: Show with text format ... "
if $BINARY show $EXAMPLES_DIR/task-definition.json -o text 2>&1 | grep -q "Family: my-app"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test show with YAML output
echo -n "Test: Show with YAML format ... "
if $BINARY show $EXAMPLES_DIR/task-definition.json -o yaml 2>&1 | grep -q "family:"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test show with JSONC
echo -n "Test: Show JSONC file ... "
if $BINARY show $EXAMPLES_DIR/task-definition.jsonc 2>&1 | grep -q "my-app"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test show with container mode
echo -n "Test: Show container definitions ... "
if $BINARY -m container show $EXAMPLES_DIR/container-definitions.json -o text 2>&1 | grep -q "web:"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

echo ""
echo "=== Running shift command tests ==="

# Test shift all containers
echo -n "Test: Shift all containers ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json --tag v2.0.0 2>&1 | grep -q "v2.0.0"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test shift with container filter
echo -n "Test: Shift specific container ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json --container web --tag v3.0.0 2>&1 | grep -q "my-app:v3.0.0"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test shift with image filter
echo -n "Test: Shift by image filter ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json --image nginx --tag stable 2>&1 | grep -q "nginx:stable"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test shift with YAML output
echo -n "Test: Shift with YAML output ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json --tag v4.0 -o yaml 2>&1 | grep -q "image:"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test shift from stdin
echo -n "Test: Shift from stdin ... "
if cat $EXAMPLES_DIR/task-definition.json | $BINARY shift --tag v5.0 2>&1 | grep -q "v5.0"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test container mode shift
echo -n "Test: Shift container definitions ... "
if $BINARY -m container shift $EXAMPLES_DIR/container-definitions.json --tag latest 2>&1 | grep -q "latest"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

echo ""
echo "=== Running overwrite option tests ==="

# Test overwrite option
TEST_FILE=$(mktemp)
cp "$EXAMPLES_DIR/task-definition.json" "$TEST_FILE"

echo -n "Test: Shift with overwrite ... "
$BINARY shift "$TEST_FILE" --tag v6.0.0 --overwrite > /dev/null 2>&1
if grep -q "v6.0.0" "$TEST_FILE"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test overwrite with YAML
cp "$EXAMPLES_DIR/task-definition.json" "$TEST_FILE"
echo -n "Test: Shift with overwrite and YAML ... "
$BINARY shift "$TEST_FILE" --tag v7.0.0 -o yaml -w > /dev/null 2>&1
if grep -q "v7.0.0" "$TEST_FILE"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test that file is NOT modified without overwrite
cp "$EXAMPLES_DIR/task-definition.json" "$TEST_FILE"
ORIGINAL_MD5=$(md5sum "$TEST_FILE" | awk '{print $1}')
$BINARY shift "$TEST_FILE" --tag v8.0.0 > /dev/null 2>&1
CURRENT_MD5=$(md5sum "$TEST_FILE" | awk '{print $1}')

echo -n "Test: File NOT modified without overwrite ... "
if [ "$ORIGINAL_MD5" = "$CURRENT_MD5" ]; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test container mode with overwrite
TEST_CONTAINER_FILE=$(mktemp)
cp "$EXAMPLES_DIR/container-definitions.json" "$TEST_CONTAINER_FILE"

echo -n "Test: Container mode with overwrite ... "
$BINARY -m container shift "$TEST_CONTAINER_FILE" --tag v9.0.0 -w > /dev/null 2>&1
if grep -q "v9.0.0" "$TEST_CONTAINER_FILE"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

rm -f "$TEST_FILE" "$TEST_CONTAINER_FILE"

echo ""
echo "=== Running error handling tests ==="

# Test missing tag
echo -n "Test: Shift without tag (should fail) ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json 2>&1 | grep -q "Error"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test invalid container
echo -n "Test: Shift with invalid container (should fail) ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json --container notfound --tag v1.0 2>&1 | grep -q "not found"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test invalid image
echo -n "Test: Shift with invalid image (should fail) ... "
if $BINARY shift $EXAMPLES_DIR/task-definition.json --image notfound --tag v1.0 2>&1 | grep -q "not found"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

# Test single object in container mode (should fail)
echo '{"name": "test"}' > /tmp/single_object.json
echo -n "Test: Single object in container mode (should fail) ... "
if $BINARY -m container show /tmp/single_object.json 2>&1 | grep -q "must be an array"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi
rm /tmp/single_object.json

# Test invalid output format
echo -n "Test: Invalid output format (should fail) ... "
if $BINARY show $EXAMPLES_DIR/task-definition.json -o invalid 2>&1 | grep -q "invalid"; then
    echo -e "${GREEN}PASS${NC}"
    ((passed++))
else
    echo -e "${RED}FAIL${NC}"
    ((failed++))
fi

echo ""
echo "=== Summary ==="
echo -e "Passed: ${GREEN}$passed${NC}"
echo -e "Failed: ${RED}$failed${NC}"
echo "Total: $((passed + failed))"

if [ "$failed" -gt 0 ]; then
    exit 1
fi

exit 0
