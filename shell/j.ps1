function j {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$Args
    )

    $env:EXE_ROOT = "__EXE_ROOT__"
    $dir = & "__BIN_PATH__" @Args
    if ($LASTEXITCODE -ne 0) {
        return
    }
    Set-Location $dir
}
