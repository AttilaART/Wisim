# Wisim

## Overview

A (WIP) competitive Company management game developed by Attila Pinter.

Bear in mind that most of the information in this document are subject to change.

### Game Loop

#### Early game

Each player is put in charge of an existing small company producing ~a (to be defined)
product~ coffee machines. The players are given a few in game steps to get to grips with how to
game works.

#### Mid game

The players compete against each other to make the most money while having to
bear events like recessions and shortages (etc.).

### Rough Progress

- [x] Core simulation
- [x] Core GUI
- [x] Core gameplay loop
- [x] Core online
- [x] Balancing
- [x] Polish
- [ ] Deeper / advanced features

## Running the beta

Download the latest release on the "releases" page (see the sidebar on the right).

- For Windows download the file ending in ".exe".
- For older Macs (Intel) download the file ending in "darwin-amd64"
- For newer Macs (Apple Silicon) download the file ending in "darwin-arm64"
- For Linux download the file ending in "linux-amd64"

Run the file. On MacOS and Linux you may have to make the file executable by running

```
chmod +x wisim-darwin-arn64
```

In the terminal (replacing "wisim-darwin-arn64" with the name of the file you downloaded)

## Building from source

### Step 0

Make sure you have the following programs installed:

- Go (programming language)
- Node (or a compatible JS runtime)
- Bash (or a compatible shell)

### Step 1

Clone the Github repo and navigate into it.

```Bash
git clone "https://github.com/AttilaART/Wisim"
cd Wisim
```

### Step 2

Download the Node JS runtime for each version from [their website](https://nodejs.org/en/download) as standalone binary.

- MacOS ARM64
- MacOS x64
- Windows x64
- Linux x64

### Step 3

Extract each archive and rename and move the executables to the following paths in the repo:

- Windows: build/wisim-nouveau/static/node-windows-amd64.exe
- MacOS ARM64 (Apple Silicon): build/wisim-nouveau/static/node-darwin-arm64
- MacOS x86 (Intel): build/wisim-nouveau/static/node-darwin-amd64
- Linux: build/wisim-nouveau/static/node-linux-amd64

### Step 4

Navigate to directory folder and run build script.

```
./build.sh
```
