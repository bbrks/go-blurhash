# go-blurhash [![Build Status](https://travis-ci.org/bbrks/go-blurhash.svg)](https://travis-ci.org/bbrks/go-blurhash) [![codecov](https://codecov.io/gh/bbrks/go-blurhash/branch/master/graph/badge.svg)](https://codecov.io/gh/bbrks/go-blurhash) [![GoDoc](https://godoc.org/github.com/bbrks/go-blurhash?status.svg)](https://godoc.org/github.com/bbrks/go-blurhash) [![Go Report Card](https://goreportcard.com/badge/github.com/bbrks/go-blurhash)](https://goreportcard.com/report/github.com/bbrks/go-blurhash) [![GitHub tag](https://img.shields.io/github/tag/bbrks/go-blurhash.svg)](https://github.com/bbrks/go-blurhash/releases) [![license](https://img.shields.io/github/license/bbrks/go-blurhash.svg)](https://github.com/bbrks/go-blurhash/blob/master/LICENSE)

A pure Go implementation of Blurhash. Right now, almost a straight up port of the [C](https://github.com/Gargron/blurhash) and [TypeScript](https://github.com/Gargron/blurhash.js) versions, slightly adapted to Go.

Blurhash is an algorithm that encodes an image into a short (~20-30 byte) ASCII string. When you decode the string back into an image, you get a gradient of colors that represent the original image. This can be useful for scenarios where you want an image placeholder before loading, or even to censor the contents of an image [a la Mastodon](https://blog.joinmastodon.org/2019/05/improving-support-for-adult-content-on-mastodon/).

Blurhash is written by [Dag Ågren](https://github.com/DagAgren).

## Contributing

Issues, feature requests or improvements welcome!

## Licence

This project is licensed under the [MIT License](LICENSE).

