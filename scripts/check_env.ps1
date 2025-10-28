# check_env.ps1 - Check development environment for quotebox project (PowerShell)

Write-Host "================================" -ForegroundColor Cyan
Write-Host "QuoteBox Environment Checker" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

$missingDeps = 0

# Function to check if a command exists
function Test-CommandExists {
    param($command)
    $null = Get-Command $command -ErrorAction SilentlyContinue
    return $?
}

# Check Go
Write-Host "Checking Go installation..." -ForegroundColor Yellow
if (Test-CommandExists "go") {
    $goVersion = go version
    Write-Host "✓ Go is installed: $goVersion" -ForegroundColor Green
} else {
    Write-Host "✗ Go is NOT installed" -ForegroundColor Red
    Write-Host ""
    Write-Host "  Installation options:" -ForegroundColor White
    Write-Host "  1. Download installer from: https://go.dev/dl/" -ForegroundColor White
    Write-Host "  2. Use winget: winget install -e --id GoLang.Go" -ForegroundColor White
    Write-Host "  3. Use Chocolatey: choco install golang" -ForegroundColor White
    Write-Host "  4. Or use WSL2 (recommended for Docker):" -ForegroundColor White
    Write-Host "     - Run: wsl --install" -ForegroundColor White
    Write-Host "     - Then in Ubuntu: sudo apt update; sudo apt install -y golang-go" -ForegroundColor White
    Write-Host ""
    $missingDeps++
}

# Check Docker
Write-Host ""
Write-Host "Checking Docker installation..." -ForegroundColor Yellow
if (Test-CommandExists "docker") {
    $dockerVersion = docker --version
    Write-Host "✓ Docker is installed: $dockerVersion" -ForegroundColor Green
    
    # Check if Docker daemon is running
    try {
        $null = docker ps 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ Docker daemon is running" -ForegroundColor Green
        } else {
            Write-Host "⚠ Docker is installed but daemon is not running" -ForegroundColor Yellow
            Write-Host "  Please start Docker Desktop" -ForegroundColor White
        }
    } catch {
        Write-Host "⚠ Docker is installed but daemon is not running" -ForegroundColor Yellow
        Write-Host "  Please start Docker Desktop" -ForegroundColor White
    }
} else {
    Write-Host "✗ Docker is NOT installed" -ForegroundColor Red
    Write-Host ""
    Write-Host "  Installation instructions:" -ForegroundColor White
    Write-Host "  1. Download Docker Desktop for Windows:" -ForegroundColor White
    Write-Host "     https://docs.docker.com/desktop/install/windows-install/" -ForegroundColor White
    Write-Host "  2. Follow the installation wizard" -ForegroundColor White
    Write-Host "  3. Restart your computer after installation" -ForegroundColor White
    Write-Host ""
    Write-Host "  Note: Docker Desktop requires Windows 10/11 Pro, Enterprise, or Education" -ForegroundColor White
    Write-Host "        with Hyper-V enabled, OR WSL2" -ForegroundColor White
    Write-Host ""
    $missingDeps++
}

# Check Docker Compose
Write-Host ""
Write-Host "Checking Docker Compose..." -ForegroundColor Yellow
if (Test-CommandExists "docker-compose") {
    $composeVersion = docker-compose --version
    Write-Host "✓ Docker Compose is available: $composeVersion" -ForegroundColor Green
} elseif (Test-CommandExists "docker") {
    try {
        $composeVersion = docker compose version 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ Docker Compose is available: $composeVersion" -ForegroundColor Green
        } else {
            Write-Host "✗ Docker Compose is NOT available" -ForegroundColor Red
            Write-Host "  Docker Compose is usually included with Docker Desktop" -ForegroundColor White
            $missingDeps++
        }
    } catch {
        Write-Host "✗ Docker Compose is NOT available" -ForegroundColor Red
        Write-Host "  Docker Compose is usually included with Docker Desktop" -ForegroundColor White
        $missingDeps++
    }
} else {
    Write-Host "✗ Docker Compose is NOT available" -ForegroundColor Red
    $missingDeps++
}

# Check Git
Write-Host ""
Write-Host "Checking Git installation..." -ForegroundColor Yellow
if (Test-CommandExists "git") {
    $gitVersion = git --version
    Write-Host "✓ Git is installed: $gitVersion" -ForegroundColor Green
} else {
    Write-Host "✗ Git is NOT installed" -ForegroundColor Red
    Write-Host "  Download from: https://git-scm.com/download/win" -ForegroundColor White
    Write-Host "  Or use winget: winget install -e --id Git.Git" -ForegroundColor White
    $missingDeps++
}

# Check for .env file
Write-Host ""
Write-Host "Checking for .env file..." -ForegroundColor Yellow
if (Test-Path .env) {
    Write-Host "✓ .env file exists" -ForegroundColor Green
    
    # Check for required variables
    $envContent = Get-Content .env -Raw
    if ($envContent -match "OPENROUTER_API_KEY=") {
        if ($envContent -match "OPENROUTER_API_KEY=your_openrouter_api_key_here") {
            Write-Host "⚠ OPENROUTER_API_KEY is set to placeholder value" -ForegroundColor Yellow
            Write-Host "  Please update .env with your actual OpenRouter API key" -ForegroundColor White
        } else {
            Write-Host "✓ OPENROUTER_API_KEY is configured" -ForegroundColor Green
        }
    } else {
        Write-Host "⚠ OPENROUTER_API_KEY not found in .env" -ForegroundColor Yellow
    }
} else {
    Write-Host "⚠ .env file not found" -ForegroundColor Yellow
    Write-Host "  Copy .env.example to .env and configure your settings:" -ForegroundColor White
    Write-Host "  Copy-Item .env.example .env" -ForegroundColor White
}

# Summary
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
if ($missingDeps -eq 0) {
    Write-Host "✓ All required dependencies are installed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor White
    Write-Host "1. Make sure .env file is configured with your OPENROUTER_API_KEY" -ForegroundColor White
    Write-Host "2. Run: docker compose up --build" -ForegroundColor White
    Write-Host "3. Access the app at: http://localhost:8080" -ForegroundColor White
    Write-Host "4. Access Grafana at: http://localhost:3000 (admin/admin)" -ForegroundColor White
    Write-Host "5. Access Prometheus at: http://localhost:9090" -ForegroundColor White
} else {
    Write-Host "✗ Some dependencies are missing. Please install them and run this script again." -ForegroundColor Red
    Write-Host ""
    Write-Host "For Windows users, we recommend using WSL2 for the best experience:" -ForegroundColor Yellow
    Write-Host "1. Run: wsl --install" -ForegroundColor White
    Write-Host "2. Restart your computer" -ForegroundColor White
    Write-Host "3. Open Ubuntu from Start Menu" -ForegroundColor White
    Write-Host "4. Install dependencies in WSL2 using the bash script" -ForegroundColor White
    exit 1
}
Write-Host "================================" -ForegroundColor Cyan
