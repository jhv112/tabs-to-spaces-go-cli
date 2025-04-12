<#
    .SYNOPSIS
        Runs the program's tests and builds the program;
        runs the program on its own source directory.

    .NOTES
        !!!WARNING!!!
        THE PROGRAM IS KNOWN, TO TRASH THE LOCAL `.git` DIRECTORY, IF:

        - THE FILE FILTER IS ABSENT;
        - THE FILE FILTER IS MALFORMED;
        - THE FILE FILTER IS POORLY THOUGHT-OUT;
        - AND/OR THE FILE FILTER IS POORLY IMPLEMENTED.

        IN SUCH AN EVENT, IT IS POSSIBLE BUT UNLIKELY, THAT UNCOMMITED CHANGES
        AND UNPUSHED COMMITS ARE SALVAGEABLE*; AND YOUR BEST BET IS, TO CLONE
        A FRESH COPY OF THE REPOSITORY FROM THE REMOTE REPOSITORY.

        *I THANK THE BLACK ARTS OF INITIALISING AND PUSHING THIS REPOSITORY ONTO
        GITHUB FOR SONEHOW SALVAGING AND THEREFORE SAVING THIS VERY REPOSITORY'S
        VERSION HISTORY.
        PRAISE BE THE OMNISSIAH!
#>

go test ./...

if (-not $?) {
    return
}

go build -o tabstospaces.exe

if (-not $?) {
    return
}

.\tabstospaces.exe --filter=.*\.go$ --tabsize=4
