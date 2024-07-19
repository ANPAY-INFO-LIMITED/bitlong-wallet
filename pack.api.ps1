$start_time = Get-Date
$currentPath = Get-Location
$currentFolderName = Split-Path -Path $currentPath -Leaf
$specificString = "wallet"
if ($currentFolderName -eq $specificString) {
    Write-Host "gomobile is in progress, please wait..."
    Set-Location api
    gomobile bind -target android
    $packtime = Get-Date -Format "yyyy-MM-dd_HHmmss"
    New-Item -Path ../ -Name "$packtime" -ItemType "directory"
    Copy-Item ./api-sources.jar ../$packtime/
    Copy-Item ./api.aar ../$packtime/
    Set-Location ../$packtime/
    $end_time = Get-Date
    $time_taken = $end_time - $start_time
    Write-Host "Time cost: $($time_taken.TotalSeconds) seconds."
    pause
} else {
	Write-Output "Wrong current directory, please run script in wallet."
    pause
}
