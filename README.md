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

#### Endgame

As the competing companies grow, they are given the ability to buy each others' stocks.
When a single company owns more than 50% of another's stocks, they can commit
to a "hostile takeover", where they buy out all of the assets of their competit-
or and the player, whose company was bought, looses. When only one player is left,
they win.

### Rough Progress

- [x] Core simulation
- [x] Core GUI
- [x] Core gameplay loop
- [x] Core online
- [ ] Balancing
- [ ] Polish
- [ ] Deeper / advanced features

## Running the beta

Before running the program, make sure to have the following programs installed

- Node (npm)
- Go

1. Clone the Git repository onto your computer.

```zsh
git clone https://github.com/AttilaART/Wisim.git
```

2. Install dependencies

```zsh
cd wisim-nouveau
npm install
cd ..
```

3. Start the servers

```zsh
./start.sh
```

After which you should be able to navigate to [localhost:5173](http://localhost:5173/).
