# gois

[![CI](https://github.com/mzbaulhaque/gois/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/mzbaulhaque/gois/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/mzbaulhaque/gois/branch/master/graph/badge.svg)](https://codecov.io/gh/mzbaulhaque/gois)
[![Go Report Card](https://goreportcard.com/badge/github.com/mzbaulhaque/gois)](https://goreportcard.com/report/github.com/mzbaulhaque/gois)

**gois** is a command line program to search images from popular services like google, bing.

## Installation

You can download the binaries directly from the [releases](https://github.com/mzbaulhaque/gois/releases) section.

## Usage

You should be able to use it directly from terminal window.

```shell
gois --help # view available commands and flags
```

Search images:

```shell
gois bing "night sky" # using bing
gois flickr "night sky" # using flickr
gois google "night sky" # using google
gois yahoo "night sky" # using yahoo
gois yandex "night sky" # using yandex
```

Please note that **gois** currently does not support automatic downloading of images. But you can easily do that using output from **gois** and piping that to **curl**/**wget**.

```shell
gois google -c "night sky" | xargs -I url curl --progress-bar --compressed --connect-timeout 10 --retry 3 -k -L -O url
gois google -c "night sky" | wget -q --show-progress -c -nc -T 10 -t 3 -i-
```

## Supported Services

  * Bing
  * Flickr
  * Google
  * Yahoo
  * Yandex
