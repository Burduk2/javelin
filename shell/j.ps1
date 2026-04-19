function j {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$Args
    )

    $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
    $exeRoot = (Resolve-Path (Join-Path $scriptDir "..\bin")).Path
    $dir = & "$exeRoot\javelin.exe" @Args
    if ($LASTEXITCODE -ne 0) {
        return
    }
    Set-Location $dir
}
