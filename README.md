## 🧤 mitten 🧤

`mitten` is a drop-in replacement for SSH that brings internet connection
to the machines without it.
For now, works only as an HTTP/HTTPS proxy.

### Building

0. Install Go (https://go.dev, `brew install go`, `conda install go`)

1. Build `mitten`:
```shell
go install github.com/unkaktus/mitten@latest
```
2. Add `$HOME/go/bin` to your `$PATH` :
```shell
export PATH="$HOME/go/bin:$PATH"
```
You would probably want to have it permanently, so put into your shellrc.

### Example use
Normally, connection fails:
```
$ ssh supercomp
supercomp> curl -I unkaktus.art
curl: (7) Failed to connect to unkaktus.art port 80: Network is unreachable
```

With `mitten`, it doesn't:

```
$ mitten supercomp
supercomp> curl -I unkaktus.art
HTTP/1.1 308 Permanent Redirect
```