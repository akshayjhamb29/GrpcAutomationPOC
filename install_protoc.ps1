# Script to download and install Protocol Buffers compiler (protoc) on Windows
# Run this script as Administrator

# Configuration
$protocVersion = "24.4"
$installDir = "C:\protoc"
$downloadUrl = "https://github.com/protocolbuffers/protobuf/releases/download/v$protocVersion/protoc-$protocVersion-win64.zip"
$zipFile = "$env:TEMP\protoc-$protocVersion-win64.zip"

# Create installation directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    Write-Host "Creating directory $installDir..."
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Download protoc
Write-Host "Downloading Protocol Buffers compiler v$protocVersion..."
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -Uri $downloadUrl -OutFile $zipFile

# Extract the zip file
Write-Host "Extracting files to $installDir..."
Expand-Archive -Path $zipFile -DestinationPath $installDir -Force

# Add to PATH environment variable
$binPath = "$installDir\bin"
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")

if ($currentPath -notlike "*$binPath*") {
    Write-Host "Adding $binPath to system PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$binPath", "Machine")
    $env:Path = "$env:Path;$binPath"
    Write-Host "Added to PATH successfully."
} else {
    Write-Host "$binPath is already in PATH."
}

# Clean up
Remove-Item $zipFile -Force
Write-Host "Temporary files cleaned up."

# Install Go plugins for Protocol Buffers
Write-Host "Installing Go plugins for Protocol Buffers..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Verify installation
Write-Host "Verifying installation..."
try {
    $protocVersion = & protoc --version
    Write-Host "Protocol Buffers compiler installed successfully: $protocVersion"
    
    Write-Host "`nInstallation complete!"
    Write-Host "You may need to restart your command prompt or IDE to use protoc."
    Write-Host "To generate Go code from your proto file, run:"
    Write-Host "cd c:\Users\aksha\Desktop\MockingDemo"
    Write-Host "protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative service.proto"
} catch {
    Write-Host "Failed to verify installation. Please check if protoc is in your PATH."
    Write-Host "You may need to restart your command prompt or IDE."
}