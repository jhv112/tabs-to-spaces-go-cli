# Tabs to Spaces CLI Utility

```yaml
Implementation language: Go
Build tool: go build
```

I'll be honest: I didn't expect, to make this a repository;
I just wanted [a quick script, that'd do some regex on files](tabstospaces.ps1).

Also, the irony of using a language, where tabs are enforced over spaces, to
convert tabs to spaces, is not lost on me.

## Warning

The program **should** work properly; however, even when functioning properly,
this program **can** demonstrate [Uncle Ben's proverb][^1] and mangle your
files, if great care is not taken.

For a concrete account of what has once gone wrong, go read the notes section
in the [dog-fooding script](dogfood.ps1) as well as the considerations being
made in the current version of the [to-do list](todo.md).

## Usage

See [USAGE](USAGE).  
Example usage can be found in the [dog-fooding script](dogfood.ps1).

## Licence

This project is licensed under the terms of the [MIT license](LICENCE).

[^1]: [With great power comes great responsibility.](https://en.wikipedia.org/wiki/With_great_power_comes_great_responsibility)
