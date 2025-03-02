# Script to run automation tests for Service A

param (
    [string]$Environment = "stg"
)

# Validate environment
if ($Environment -ne "stg" -and $Environment -ne "qa") {
    Write-Host "Invalid environment: $Environment. Must be 'stg' or 'qa'"
    exit 1
}

Write-Host "Running Service A tests in $Environment environment"

# Install test dependencies if needed
Write-Host "Installing test dependencies..."
go get github.com/stretchr/testify/assert
go get google.golang.org/grpc/test/bufconn

# Run the tests
Write-Host "Running tests..."
cd c:\Users\aksha\Desktop\MockingDemo
go test -v ./automation -env=$Environment

Write-Host "Tests completed."