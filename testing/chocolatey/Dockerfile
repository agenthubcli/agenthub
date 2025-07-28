# escape=`
FROM mcr.microsoft.com/windows/servercore:ltsc2022

# Install Chocolatey
RUN powershell -Command `
    Set-ExecutionPolicy Bypass -Scope Process -Force; `
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; `
    iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Copy package files
COPY . /package
WORKDIR /package

# Test the chocolatey script syntax
RUN powershell -Command `
    "Write-Host 'Testing PowerShell script syntax...'; `
    try { `
        . ./tools/chocolateyInstall.ps1 -WhatIf -ErrorAction Stop; `
        Write-Host 'PowerShell script syntax is valid!' -ForegroundColor Green `
    } catch { `
        Write-Host \"PowerShell script error: \$_\" -ForegroundColor Red; `
        exit 1 `
    }"

# Build and test the package
RUN powershell -Command `
    "choco pack --trace; `
    if (\$LASTEXITCODE -ne 0) { exit \$LASTEXITCODE }; `
    Write-Host 'Package built successfully!' -ForegroundColor Green"

CMD ["powershell", "-Command", "Write-Host 'Chocolatey package test completed successfully!' -ForegroundColor Green"]
