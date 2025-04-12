# test, build, run on own directory
# !!!IS KNOWN, TO TRASH THE `.git` DIRECTORY, IF FILE FILTER DOESN'T WORK!!!

go test ./...

if (-not $?) {
    return
}

go build -o tabstospaces.exe

if (-not $?) {
    return
}

.\tabstospaces.exe --filter=.*\.go$ --tabsize=4
