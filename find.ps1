$agentIds = @(1079, 2378, 2184, 3379, 4794)
$startDate = @(2010, 2019)
Write-Host $agentIds
Write-Host $PSScriptRoot

Foreach ($agent in $agentIds)
{
    Foreach ($date in $startDate)
    {
        for ($j=0; $j -lt 3; $j++){
            Start-Process -NoNewWindow -FilePath ".\cmd\find_mongo\find_mongo.exe" -ArgumentList "$agent $date-01-01 2020-12-31" -Wait
            # Start-Process -NoNewWindow -FilePath ".\cmd\find_mssql\find_mssql.exe" -ArgumentList "$agent $date-01-01 2020-12-31" -Wait
            # Start-Process -NoNewWindow -FilePath ".\cmd\find_neo4j\find_neo4j.exe" -ArgumentList "$agent $date-01-01 2020-12-31" -Wait
        }
    }
}