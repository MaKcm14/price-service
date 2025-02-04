try {
    ## Starting the auxiliary services.
    # Starting the by-pass-service.
    $by_pass_service = Start-Process -PassThru -FilePath "python" -ArgumentList "../../tools/repository/api/by_pass_service/server.py"

    ## Starting the main service.
    $price_service = Start-Process -PassThru -FilePath "go" -ArgumentList "run", "../../cmd/app/main.go"

    $env:BY_PASS_SERVICE_ID=$by_pass_service.Id
    $env:PRICE_SERVICE_ID=$price_service.Id

} catch {
    Write-Output "ERROR: The service wasn't started correctly. Please note the errors and try again."
}
