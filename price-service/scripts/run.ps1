## Starting the auxiliary services.
# Starting the by-pass-service
Start-Process -FilePath "python" -ArgumentList "../tools/repository/api/by_pass_service/server.py"

## Starting the main service
Start-Process -FilePath "go" -ArgumentList "run", "../cmd/app/main.go"

Write-Output "The service was configured and started"
