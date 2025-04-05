# tabs to spaces

Get-ChildItem -Include '*.H', '*.CPP' -Recurse | % {
    $content = Get-Content -Path $_ -Raw -Encoding UTF8;

    while (($content | ? { $_ -Match '\t' })) {
        $content = $content `
            -Replace '((?:\n|^)(?:[^\t\n]{4})*[^\t\n]{0})\t', '$1    ' `
            -Replace '((?:\n|^)(?:[^\t\n]{4})*[^\t\n]{1})\t', '$1   ' `
            -Replace '((?:\n|^)(?:[^\t\n]{4})*[^\t\n]{2})\t', '$1  ' `
            -Replace '((?:\n|^)(?:[^\t\n]{4})*[^\t\n]{3})\t', '$1 '
    }

    # from: https://stackoverflow.com/questions/5596982/a/32951824
    [IO.File]::WriteAllLines($_, $content)

    # # writes UTF8 with BOM
    #  $content | Out-File $_ -Encoding utf8 -NoNewline
}
