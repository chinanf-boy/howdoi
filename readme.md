## howdoi

Copy from [howdoi:py](https://github.com/gleitz/howdoi), But no done yet

### Install

```
go get -v -u github.com/chinanf-boy/howdoi
```

### TODO

- [x] base Feature, Get the data
- [x] Proxy can with [Socks5](./howdoi/client.go)
- [x] how many answers you want
- [ ] colorful Code text with shell env
- [ ] cache Result ?
- [ ] test file

> **Notes:** , cli-name Over the Python version cli

### Uasge

same as py:howdoi , but lit diff

``` bash
$ howdoi -q "format date bash"
```

#### Tips

> About the `ENV`

ENV | Desc | Default
---------|----------|---------
| **HOWDOI_DISABLE_SSL** | change `https://` => `http://` | `nil`
| **HOWDOI_URL** | search engine with the ask website  | `stackoverflow.com`
| **HOWDOI_SEARCH_ENGINE** | search engine{bing\|google} | `bing`

### Cli

``` js
usage: howdoi [-h|--help] [-c|--color] [-v|--version] [-n|--num <integer>]
              -q|--query "<value>" [-q|--query "<value>" ...]

              cli to Ask the question

Arguments:

  -h  --help     Print help information
  -c  --color    colorful Output. Default: false
  -v  --version  version
  -n  --num      how many answer. Default: 1
  -q  --query    query what
```

### Why rewrite

1. proxy,[some issue](https://github.com/chinanf-boy/howdoi/issues/1) with `socks`
2. fast


