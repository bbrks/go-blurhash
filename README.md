# go-blurhash [![Go Reference](https://pkg.go.dev/badge/github.com/bbrks/go-blurhash.svg)](https://pkg.go.dev/github.com/bbrks/go-blurhash) [![GitHub tag](https://img.shields.io/github/tag/bbrks/go-blurhash.svg)](https://github.com/bbrks/go-blurhash/releases) [![license](https://img.shields.io/github/license/bbrks/go-blurhash.svg)](https://github.com/bbrks/go-blurhash/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/bbrks/go-blurhash)](https://goreportcard.com/report/github.com/bbrks/go-blurhash) [![codecov](https://codecov.io/gh/bbrks/go-blurhash/branch/master/graph/badge.svg)](https://codecov.io/gh/bbrks/go-blurhash)

A pure Go implementation of [Blurhash](https://github.com/woltapp/blurhash). The API is stable, however the hashing function in either direction may not be.

![Blurhash Demo](https://i.imgur.com/9qxOXJW.png)

Blurhash is an algorithm written by [Dag Ågren](https://github.com/DagAgren) for [Wolt (woltapp/blurhash)](https://github.com/woltapp/blurhash) that encodes an image into a short (~20-30 byte) ASCII string. When you decode the string back into an image, you get a gradient of colors that represent the original image. This can be useful for scenarios where you want an image placeholder before loading, or even to censor the contents of an image [a la Mastodon](https://blog.joinmastodon.org/2019/05/improving-support-for-adult-content-on-mastodon/).

Under the covers, this library is almost a straight port of the [C version](https://github.com/woltapp/blurhash/tree/master/C), which is known to encode images slightly differently than the TypeScript implementation.

## Contributing

Issues, feature requests or improvements welcome!

## Licence

This project is licensed under the [MIT License](LICENSE).
