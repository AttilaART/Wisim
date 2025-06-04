# Wisim

## Overview

A (WIP) competitive Company management game developed by Attila Pinter.

Bear in mind that most of the information in this document are subject to change.

### Game Loop

#### Early game

Each player is put in charge of an existing small company producing a (to be defined)
product. The players are given a few in game steps to get to grips with how to
game works.

#### Mid game

The players compete against each other to make the most money while having to
bear events like recessions and shortages (etc.).

#### Endgame

As the competing companies grow, they are given the ability to buy each others' stocks.
When a single company owns more than 50% of another's stocks, they can commit
to a "hostile takeover", where they buy out all of the assets of their competit-
or and the player, whose company was bought, looses. When only one player is left,
they win.

### Rough Progress

- [x] Core simulation
- [ ] Core GUI
- [ ] Core gameplay loop
- [ ] Core online
- [ ] Polish

## Running the beta

Before running the program, make sure to have the following programs installed

- Node (npm)
- Go
- Wails CLI

1. Clone the Git repository onto your computer

```zsh
git clone https://github.com/AttilaART/Wisim.git
```

2.

```zsh
cd Wisim/Wisim
wails dev
```

After which the program will start by itself.

If none of this works you can try running the compiled binaries to be found in the
`Wisim/Wisim/build/bin` directory.

Make sure to run the correct one for your device.

MacOS: `WiSim.app`
Windows: `WiSim.exe`
Linux: `WiSim`
