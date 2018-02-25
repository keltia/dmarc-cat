# README.md

## Status

*stable*
[![Build Status](https://secure.travis-ci.org/keltia/dmarc-cat.png)](http://travis-ci.org/keltia/dmarc-cat)

## Installation

As with many Go utilities, a simple

    go get github.com/keltia/dmarc-cat

is enough to fetch, build and install.

## Dependencies

Aside from the standard library, I use `github.com/intel/tfortools` to generate tables.

    go get -u github.com/intel/tfortools

## Usage

SYNOPSIS
```
darc-cat -hvD <zipfile|xmlfile>

Example:

$ dmarc-cat /tmp/yahoo.com\!keltia.net\!1518912000\!1518998399.xml

Reporting by: Yahoo! Inc. — postmaster@dmarc.yahoo.com
From 2018-02-18 01:00:00 +0100 CET to 2018-02-19 00:59:59 +0100 CET

Domain: keltia.net
Policy: p=none; dkim=r; spf=r

Reports(1):
IP            Count   From       RFrom      RDKIM   RSPF
88.191.250.24 1       keltia.net keltia.net neutral pass
```

## Tests

Very incomplete for now.

## License

This is released under the BSD 2-Clause license.  See `LICENSE.md`.

## Contributing

I use Git Flow for this package so please use something similar or the usual github workflow.

1. Fork it ( https://github.com/keltia/dmarc-cat/fork )
2. Checkout the develop branch (`git checkout develop`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create a new Pull Request
