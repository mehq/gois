[![CI](https://github.com/mzbaulhaque/gomage/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/mzbaulhaque/gomage/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/mzbaulhaque/gomage/branch/master/graph/badge.svg)](https://codecov.io/gh/mzbaulhaque/gomage)

**gomage** - CLI program to search and download images in bulk

## INSTALLATION

To install from source, use the following command:

```shell
make install # Installs gomage at $GOBIN or $GOPATH/bin
```

**Go (v1.16+) is required.**

## USAGE

You can use it directly from any terminal window after installation (assuming installation directory is in $PATH).

```shell
gomage [options] query # query is any valid search query, e.g. cats, dogs
```

To view available options:
```shell
gomage -help
```
