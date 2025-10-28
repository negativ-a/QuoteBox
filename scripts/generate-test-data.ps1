#!/usr/bin/env pwsh

Write-Host "`n🎯 GENERATING TEST QUOTES..." -ForegroundColor Cyan

$tags = @('joy', 'love', 'hope', 'gratitude', 'peace', 'courage', 'wisdom')
$counts = @(5, 4, 4, 3, 3, 2, 2)

for($i=0; $i -lt $tags.Length; $i++) {
    Write-Host "`nGenerating $($counts[$i]) quotes for '$($tags[$i])'..." -ForegroundColor Yellow
    
    1..$counts[$i] | ForEach-Object {
        $jsonBody = @{tag = $tags[$i]} | ConvertTo-Json -Compress
        $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/quote" -Method Post -Body $jsonBody -ContentType "application/json"
        $preview = $response.quote.Substring(0, [Math]::Min(50, $response.quote.Length))
        Write-Host "  ✓ Quote $($_): $preview..." -ForegroundColor Green
        Start-Sleep -Milliseconds 700
    }
}

$total = ($counts | Measure-Object -Sum).Sum
Write-Host "`n✅ Generated $total quotes across $($tags.Length) tags!" -ForegroundColor Green
