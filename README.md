# NWGO - CLI

a command-line tool, that is aiming to combine the flexibility of a NW.JS ([https://nwjs.io](https://nwjs.io)) UI
with the power of a Go based backend and compile everything into one Executable in the End.


* [Installation](#installation)
* [Creating a Project](#creating-a-project)
* [Command-Reference](#command-reference)
* [Roadmap](#roadmap)


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
(Only Mac ARM, Linux x86_64 are are supported at the moment)

everything this tool downloads is stored in ~/.local/state/nwgo

# Creating a Project

1.) create a new Project in a Directory of you choice.

```bash
nwgo init [ProjectDirectoryName]
```

2.) you'll be asked for a **"NWJS - Package"** name
this is the name that represents the NWJS-Frontend part of your application.

3.) then you'll be asked for as **"Go-Package"** name
this is the Go-Backend side of your application.  
an example would be `github.com/your-account/your-project-directory`

> [!info]  
> both inputs give you rules on what characters are allowed in the respective names

4.) if everything fits, you should now have a Folder, containing the following structure

```
[ProjectDirectoryName]
    ├── go.mod
    ├── goapi
    │   └── server.go
    ├── index.html
    ├── main.go
    └── package.json
```

5.) Switch into the folder that was created.

```bash
cd [ProjectDirectoryName]
```

6.) Run the Project to see if everything works

```bash
nwgo run .
```

> [!info]  
> The first time you run this command, it will download and set up the NWJS-SDK-For your Platform.
> This may take a while, depending on your Download speed.

7.) your Project will live in the `goapi` folder.  
You define all your Backend-Routes in `goapi/server.go`.  
The Entry Point for your Application is the `GET /` Route. There shoudl already
be an example created for you.

# Command-Reference

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
    -   ▢ for Window
    -   ✅ for Mac (arm)
    -   ~~for Mac (x86)~~  
        <small>(can't verify due to missing Hardware)</small>
    -   ✅ for Linux (x86)
    -   ~~for Linux (arm)~~  
        <small>(NW.JS does not seem to have a Linux-Arm64 Build)</small>

Step 2:

-   ✅ add a go based server template to the project created by `nwgo init`
-   ✅ have `nwgo` automatically compile that `go-server` project
-   ✅ have `nwjs` launch that `go server` on a free port uppon start
-   ✅ pass the port of the launched file back to `nwjs`
-   ✅ make sure the server is killed when nwjs closes
-   ▢ fix Error-Output during backend compile process

Step 3:

-   ▢ add a `nwgo build [path] [target]` command
-   ▢ download the none SDK version of NWJS
-   ▢ combine the project into a 'nw' file
-   ▢ link the `nw-file` together with the download NWJS to create one application.

Step 4

-   ▢ fix window-close-detection, when switching to a different domain
