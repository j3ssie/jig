Jig - Jaeles Input Generator
========

This is helper tool for generated inputs for [Jaeles](https://github.com/jaeles-project/jaeles) Scanner.

## Install

```shell
GO111MODULE=on go get github.com/j3ssie/jig
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

Available Output Type:
  location   --  Use Location headers as {{.BaseURL}}
```

## License

`Jaeles` is made with ♥  by [@j3ssiejjj](https://twitter.com/j3ssiejjj) and it is released under the MIT license.
