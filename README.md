# Bhojpur License - Management Engine

The Bhojpur License is used as a software Licensing Engine applied within the [Bhojpur.NET Platform](http://github.com/bhojpur/platform) ecosystem for distribted applications or service delivery. Generally, it adds standard copyright text to the source code files, digitally signs program's binary images, requests license keys for different purpose, verifies license expiry time, etc. realated to different software product or service offerings.

The [Bhojpur License](http://github.com/bhojpur/license) requires `Go` 1.17 or later.

## Installation

    go install github.com/bhojpur/license@latest

## Usage

    license [flags] pattern [pattern ...]

    -c copyright holder (defaults to "Bhojpur Consulting Private Limited, India.")
    -f custom license file (no default)
    -l license type: Apache, BSD, MIT, MPL (defaults to "apache")
    -y year (defaults to current year)
    -check check only mode: verify presence of license headers and exit with non-zero code if missing
    -ignore file patterns to ignore, for example: -ignore **/*.go -ignore vendor/**

The pattern argument can be provided multiple times, and may also refer
to single files.

The `-ignore` flag can use any pattern [supported by
doublestar](https://github.com/bmatcuk/doublestar#patterns).

## Running in a Docker Container

The simplest way to get the Bhojpur License docker image is to pull from GitHub
Container Registry:

```bash
docker pull ghcr.io/bhojpur/license:latest
```

Alternately, you can build it from source yourself:

```bash
docker build -t ghcr.io/bhojpur/license .
```

Once you have the image, you can test that it works by running:

```bash
docker run -it ghcr.io/bhojpur/license -h
```

Finally, to run it, mount the directory you want to scan to `/src` and pass the
appropriate Bhojpur License flags:

```bash
docker run -it ghcr.io/bhojpur/license -v ${PWD}:/src -c "Bhojpur Consulting Private Limited, India." *.go
```