## 🧤 mitten 🧤

`mitten` is a drop-in replacement for SSH that brings internet connection
to the machines without it, and enables easy file transfer between the local
and remote machines.
Only you can use the connection and have the access to the local files.

### Internet demo
Normally, on machines with no internet access, calls like these fail or hang:
```
$ ssh supercomp
supercomp> curl -I https://unkaktus.art
curl: (7) Failed to connect to unkaktus.art port 443: Network is unreachable
```

With `mitten`, they work through your local connection:

```
$ mitten supercomp

       ▗▟▀▀▙
      ▗▛   ▐▌
    ▗▟▘   ▗▛
▗▄▄▟▀     ▀▀▀▀▀▀▀▜▄
█  █              ▝▜▖
█  █                ▙
█  █               ▗▌
▜▄▄█▄            ▗▟▀
     ▀▀▀▀▄▄▄▄▄▄▄▀▀
   mitten mittens!

supercomp> curl -I https://unkaktus.art
HTTP/1.1 308 Permanent Redirect
```

Mitten magic!

### File transfer demo
To easily transfer files between the local machine and the remote,
use `mittenfs` command after logging in, which provides `sftp` interface:

```shell
$ mitten supercomp
supercomp> mittenfs
supercomp> sftp> get mitten.go .
supercomp> sftp> lls
mitten.go
```

Mitten magic!

### Easy installation using Mamba

Having [MambaForge installed](https://github.com/conda-forge/miniforge#install), install `mitten` package:
```shell
mamba install -c https://mamba.unkaktus.art mitten
```

### Installation using Go

0. Install Go (https://go.dev, `brew install go`, `conda install go`)

1. Build `mitten`:
```shell
go install github.com/unkaktus/mitten@latest
```
2. Add `$HOME/go/bin` to your `$PATH` :
```shell
export PATH="$HOME/go/bin:$PATH"
```
You would probably want to have it permanently, so put it into your shellrc.
