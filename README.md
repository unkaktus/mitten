## 🧤 mitten 🧤

`mitten` is a drop-in replacement for SSH that brings internet connection
to the machines without it, and enables easy file transfer between the local
and remote machines.
There is built-in authentication, so only the user inside of the mitten session
can use the connection and have the access to the local files.

A common use case for mitten is installing packages in Conda/Miniforge environments
which rely on fetching the them from the internet.

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
█  █      mitten    ▙
█  █      magic!   ▗▌
▜▄▄█▄            ▗▟▀
     ▀▀▀▀▄▄▄▄▄▄▄▀▀


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
supercomp> sftp> get local_file.dat .
supercomp> sftp> lls
local_file.dat
```

Mitten magic!

### Binary installation

1. Download and install the binary for your platform. For example, for Linux on AMD64:

```shell
curl -L -o mitten https://github.com/unkaktus/mitten/releases/latest/download/mitten-linux-amd64
mkdir -p ~/bin
mv mitten ~/bin/
chmod +x ~/bin/mitten
```

Add $HOME/bin into your $PATH into your .bashrc:
```shell
export PATH="$HOME/bin:$PATH"
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


### A note on OpenWRT

Certain apps, such as `opkg` of OpenWRT, don't support HTTP proxy authentication.
One can disable this authentication by setting `MITTEN_DISABLE_PROXY_AUTH=yes`.
Note that it is, of course, **insecure**.


### How it works

Behind the scenes, _mitten_ launches an HTTP(S) proxy, and starts an SSH session with
forwarding of the local proxy port to the remote machine. After the connection is made,
and the user logged in, the shell prompt is getting detected. Then, _mitten_ types in
a shell command which exports `http_proxy` environment variables with the address and
authentication credentials set pointing to the forwarded connection. Because of this,
only the current SSH session has the password available, so other users and other
sessions are not able to use the proxy. After all these steps are done, the _mitten_
thumbs-up banner is displayed, and the user can use any tools that automatically pick up
the HTTP proxy configuration to access internet resources.