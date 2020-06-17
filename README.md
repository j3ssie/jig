Jig - Jaeles Intput Generator
========

This is helper tool for generated output for [Jaeles](https://github.com/jaeles-project/jaeles) Scanner.
## Install

```shell
GO111MODULE=on go get github.com/jaeles-project/jaeles
```

## Usage

```shell
jig scan -u <url> -I location
jig scan -U urls.txt -I location -o jig-output.txt
```

then run Jaeles scanner with this output from jig

```shell
jaeles scan -s <signatures> -U jig-output.txt -J

jig scan -u <url> -I location | jaeles scan -s <signatures> -J

```

## License

`Jaeles` is made with ♥  by [@j3ssiejjj](https://twitter.com/j3ssiejjj) and it is released under the MIT license.
