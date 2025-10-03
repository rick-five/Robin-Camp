# 设置TLS协议版本
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# 定义电影数据
$body = @{
    title = "Test Movie"
    releaseDate = "2023-01-01T00:00:00Z"
    genre = "Action"
    distributor = "Test Distributor"
    budget = 1000000
    mpaRating = "PG-13"
} | ConvertTo-Json

# 发送POST请求创建电影
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/movies/" `
                                  -Method POST `
                                  -Headers @{Authorization = "Bearer test_token_12345"} `
                                  -ContentType "application/json" `
                                  -Body $body
    Write-Host "Response Status Code: $($response.StatusCode)"
    Write-Host "Response Content: $($response.Content)"
} catch {
    Write-Host "Error occurred:"
    Write-Host "Status Code: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Error Message: $($_.Exception.Message)"
    if ($_.ErrorDetails) {
        Write-Host "Error Details: $($_.ErrorDetails.Message)"
    }
}