## howdoi

Copy from [howdoi:py](https://github.com/gleitz/howdoi), But no done yet

### Install

```
go get -v -u github.com/chinanf-boy/howdoi
```

### TODO

- [x] base Feature, Get the data
- [x] Proxy can with [Socks5](./src/client.go#L44)
- [x] how many answers you want `-n`
- [x] **go** func with questions
- [x] colorful Code text with shell env `-c`/`-T`, [chroma](https://godoc.org/github.com/alecthomas/chroma)
- [x] cache Result `-C` [useful refs](https://github.com/chinanf-boy/howdoi/issues/3)
- [ ] ReCache Result ?
- [ ] test file

> **Notes:** , cli-name Over the Python version cli

### Uasge

same as py:howdoi , but lit diff

```bash
$ howdoi -q "format date bash"
```

#### Tips

About the `ENV`

| ENV                      | Desc                               | Default               |
| ------------------------ | ---------------------------------- | --------------------- |
| **HOWDOI_DISABLE_SSL**   | change `https://` => `http://`     | `nil`                 |
| **HOWDOI_URL**           | search engine with the ask website | `stackoverflow.com`   |
| **HOWDOI_SEARCH_ENGINE** | search engine{bing\|google}        | `bing`                |
| **HOWDOI_CACHE_DIR**     | http Response - Cached dir         | `$HOME/.howdoi-cache` |
| **-T**                   | [chroma theme](#chroma-theme)      | `pygments`            |

> **NOTE**, careful about **Cache dir**, you will miss the ever Data after you changed diff HOWDOI_CACHE_DIR。

### Ref

#### chroma theme

<details>

<summary> info </summary>

```go
[
  abap, algol, algol_nu, arduino, autumn, borland, bw, colorful, dracula, emacs, friendly, fruity, github, igor, lovelace, manni, monokai, monokailight, murphy, native, paraiso-dark, paraiso-light, pastie, perldoc, pygments, rainbow_dash, rrt, solarized-dark, solarized-dark256, solarized-light, swapoff, tango, trac, vim, vsxcode
]
```

</details>

### Cli

```js
usage: howdoi [-h|--help] [-c|--color] [-v|--version] [-n|--num <integer>]
              -q|--query "<value>" [-q|--query "<value>" ...] [-D|--debug]
              [-T|--theme "<value>"] [-C|--cache]

              cli to Ask the question

Arguments:

  -h  --help     Print help information
  -c  --color    colorful Output. Default: false
  -v  --version  version
  -n  --num      how many answer. Default: 1
  -q  --query    query what
  -D  --debug    debug *
  -T  --theme    chrome styles. Default: pygments
  -C  --cache    cache response? . Default: false
```

### Why rewrite

1. proxy,[some issue](https://github.com/chinanf-boy/howdoi/issues/1) with `socks`
2. fast
