param(
    [string[]] $functions = @()
)

$srcRoot = Join-Path $PSScriptRoot "src"
$functionsRoot = Join-Path $srcRoot "api" "functions"
$buildRoot = Join-Path $PSScriptRoot "build"

if ($functions.count -eq 0) {
    $functions = Get-ChildItem -Path $functionsRoot -Directory | Select-Object -ExpandProperty Name
}

Write-Output "Building functions: $functions"

$env:GOARCH="amd64"
$env:GOOS="linux"

Push-Location $functionsRoot
Write-Output "Generating wire_gen.go files"
Write-Output $env:GOPATH
foreach ($func in $functions) {
    <# $func is the current item #>
    Write-Output "${func}: generating wire_gen.go"
    wire "$functionsRoot/$func/"
    if ($? -eq $false) {
        Write-Error "${func}: wire failed"
        exit 1
    }

    Write-Output "${func}: building"
    go build -o "$buildRoot/$func/bootstrap" "$functionsRoot/$func/"
    if ($? -eq $false) {
        Write-Error "${func}: go build failed"
        exit 1
    }

    Write-Output "${func}: done"
}

Pop-Location
Remove-Item Env:GOARCH
Remove-Item Env:GOOS

