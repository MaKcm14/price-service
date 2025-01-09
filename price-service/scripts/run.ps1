## Starting the auxiliary services.
# Starting the by-pass-service
$by_pass_service = Start-Process -PassThru -FilePath "python" -ArgumentList "../tools/repository/api/by_pass_service/server.py"

## Starting the main service
$price_service = Start-Process -PassThru -FilePath "go" -ArgumentList "run", "../cmd/app/main.go"

Write-Output "The service was configured and started"
