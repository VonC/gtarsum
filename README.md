# gtarsum

Compute hash for a tar file

## Usage

### Normal: output hash

```bash
gtarsum <afile.tar>
a10840209cb1c93c6bb85a34e969cf7eaaf43128b477f0f900cac49b551d26bd
```

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
- sort the list of file names
- compute a hash from the concatenation of the filnames hashes

So if a different tar file has the same files but in a different order, or different date/owner, the global hash can still be the same as an previous archive.