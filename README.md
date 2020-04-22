# nicehttp

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/lithdew/nicehttp)
[![Discord Chat](https://img.shields.io/discord/697002823123992617)](https://discord.gg/HZEbkeQ)

Package nicehttp contains helper utilities for downloading files/making requests with [valyala/fasthttp](https://github.com/valyala/fasthttp).

- Download a file from a URL serially/in chunks with multiple workers in parallel, should the URL allow it.
- Download contents of a URL and write its contents to a `io.Writer`.
- Query the headers of a URL using a HTTP head request.
- Follow redirects provisioned by a URL.