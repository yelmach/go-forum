# Set environment variables
$env:CGO_ENABLED=1

# Verify GCC installation
gcc --version

# If you need to permanently add MinGW to your PATH (run as administrator)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "Machine") + ";C:\mingw64\bin",
    "Machine"
)

# Rebuild your application
go mod tidy
go build