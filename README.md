# cazelle

Generate `BUILD.bazel` for `rules_cue`

## Install

```sh
go get github.com/abcue/cazelle
```

## Use

```sh
go run github.com/abcue/cazelle -help
Usage of .../cazelle
cazelle:
  -template string
        Path to template file
```

## Develop

Go code is generated from [this search](https://www.perplexity.ai/search/get-all-imported-cuelang-packa-TnTtLs06Q5CM_bGpcfHMqA) powered by DeepSeek R1@Perplecity.

```sh
go run .
go test -v
```
