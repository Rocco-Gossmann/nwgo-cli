# NWGO - CLI

a command-line tool, that is aiming to combine the flexibility of a NW.JS ([https://nwjs.io](https://nwjs.io)) UI
with the power of a Go based backend and compile everything into one Executable in the End.

# Installation

make sure you use **GO 1.22.0 or later**

clone this Repo:

```bash
git clone https://github.com/rocco-gossmann/nwgo-cli
```

switch into the folder, that was just created

```bash
cd nwgo-cli
```

Run the Install command.
Sadly Gos `install` does not support changing the target binary name yet, so
nwgo-cli comes with its own `install` command to compensate for that.

run this commands

```bash
go mod tidy 
go run . install
```

### Do I need to install NW.JS ?

No, you don't. The first time you use the `nwgo run` command, this tool will
automatically download the fiting NW.JS Files for your system.

everything this tool downloads is stored in ~/.local/state/nwgo

# Usage

Once you installed it, you can access the

```bash
nwgo
```

command.

typing it by itself gives you a help, for what it can do.

| Command            | Description                                                                                                                 |
| ------------------ | --------------------------------------------------------------------------------------------------------------------------- |
| `nwgo install`     | Installs the `nwjs` binary into your `$GOPATH/bin`                                                                          |
| `nwgo init [path]` | creates a new NWJS - Project at the given `[path]`                                                                          |
| `nwgo run  [path]` | launchens the given NWJS - Project at the given `[path]`                                                                    |
| `nwgo uninstall`   | Removes everything, that has been downloaded or installed via this tool<br><small>Your Projects will stay untouched</small> |

# Roadmap

Step 1:

-   ✅`nwgo` install
-   ✅ auto-download NWJS - Binaries on first run
    -   ✅ for Mac (arm)
    -   ▢ for Window
    -   ▢ for Linux (x86)
    -   ▢ for Linux (arm)

Step 2:

-   ▢ add a go based server template to the project created by `nwgo init`
-   ▢ have `nwgo` automatically compile that `go-server` project
-   ▢ have `nwjs` launch that `go server` on a free port uppon start
-   ▢ pass the port of the launched file back to `nwjs`
-   ▢ make sure the server is killed when nwjs closes
-   ▢ make sure the server is killed when nwjs closes

Step 3:

-   ▢ add a `nwgo build [path] [target]` command
-   ▢ download the none SDK version of NWJS
-   ▢ combine the project into a 'nw' file
-   ▢ link the `nw-file` together with the download NWJS to create one application.
