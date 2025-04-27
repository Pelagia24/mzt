docker-compose -f docker-compose.test.yml up -d

Write-Host "Waiting for database to be ready..."
Start-Sleep -Seconds 5

go test ./internal/repository -v

docker-compose -f docker-compose.test.yml down -v

go test ./internal/service -v