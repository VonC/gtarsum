[![Go Report Card](https://goreportcard.com/badge/Vonc/gtarsum)](https://goreportcard.com/report/Vonc/gtarsum)

# gtarsum

Compute hash for one tar file (or several tar files concurrently)

## Installation

### go get

```go
go get github.com/VonC/gtarsum
```

### releases

Artifacts from the [latest release](https://github.com/VonC/gtarsum/releases/latest) for Mac, Linux and Windows (64 bits only)


## Usage

### Normal: output hash

```bash
gtarsum <afile.tar>

a10840209cb1c93c6bb85a34e969cf7eaaf43128b477f0f900cac49b551d26bd
```

### Multiple archives: hash comparison

```bash
gtarsum <afile.tar> <afile2.tar> ...

echo $?
# or
echo %ERRORLEVEL%
```

- compute hash for each archive concurrently, not sequentially
- exit 0 if all archives are identical
- exit 1 if one or several archives differ from the first


### Verbose: with progress bar, copy hash to a file

- If the environment variable `progress` is set (to anything), it will display a progress bar for each files read in the tar file.

- If the environment variable `progress` is set (to a value ending with '`.hash`'), it will copy the result in the file '`xxx.hash`'.

```bash
progress=1.hash gtarsum <afile.tar>

File 'ex.tar' (73): 100% [============================================================================]
File 'ex.tar' hash='a10840209cb1c93c6bb85a34e969cf7eaaf43128b477f0f900cac49b551d26bd'

cat 1.hash
a10840209cb1c93c6bb85a34e969cf7eaaf43128b477f0f900cac49b551d26bd
```

### Verbose: with progress bar, for multiple archives

```bash
progress=2.hash gtarsum <afile.tar> gtarsum <afile2.tar> ...

File 'ex.tar' (73): 100% [============================================================================]
File 'ex2.tar' (132): 100% [==========================================================================]
File 'ex.tar' hash='a10840209cb1c93c6bb85a34e969cf7eaaf43128b477f0f900cac49b551d26bd'
File 'ex2.tar' hash '543a12be3e27d85e94cdfac3eae186cd7d54d4994ccd3db0f96a8077578a6bed' differs

cat 2.hash1
a10840209cb1c93c6bb85a34e969cf7eaaf43128b477f0f900cac49b551d26bd
cat 2.hash2
543a12be3e27d85e94cdfac3eae186cd7d54d4994ccd3db0f96a8077578a6bed
```

- compute sha256 for each archive, concurrently (not sequentially)
- display hash for each archive
- same rule as before: if the environment variable `progress` ends with `.hash`, the hashes are accessible in `xxx.hash1`/`xxx.hash2`/`...` files.

### Version

Use `-v`, or `--version` or `version`

```bash
gtarsum -v
Git Tag   : v0.0.1-3-gfc947e7
Build User: VonC
Version   : v0.0.1
BuildDate : 20200904-105250
```

## Tar hash

- Compute hash for each files found in the tar
- Sort the list of file names part of the archive
- Compute a hash from the concatenation of the filnames hashes

So if a different tar file has the same files but in a different order, or different date/owner, the global hash can still be the same as an previous archive.