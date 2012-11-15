# Go Header Update

Go Header Update (ghu) replaces the content between `/* */` at the beginning of
a file (A) with the content of another file (B). The `/* */` markers must start
at the beginning of file (A). Whitespace characters are ignored. If no markers
are found in (A), the content from (B) is added at the beginning of (A).

ghu is written in Go.

## Options

	-path    The path to process.
	-hfile   The file with the header content.
	-ftype   The file type pattern (suffix) to process.

## Newlines

The file with the header content is expected to have no trailing newlines,
except the required ones. Beware of vim's eol feature.

## Use Case

ghu could be useful to replace header information in source code files.

	./ghu -path <source_dir> -hfile <header_tpl> -ftype *.cpp
