# DrakkarRoller

A small Go utility that connects to a local game server over TCP and repeatedly
re-rolls a new character's stats until they meet the minimums configured in
`config.json`, then accepts the roll.

## Configuration

Edit `config.json`:

```json
{
  "CharacterName": "Smacks",
  "CharacterIsMale": true,
  "MinimumStr": 17,
  "MinimumInt": 17,
  "MinimumWis": 1,
  "MinimumWil": 17,
  "MinimumCon": 1,
  "MinimumAgi": 17,
  "MinimumCha": 17,
  "MinimumLuck": 17
}
```

Each `Minimum*` field is the lowest acceptable value for that stat; rolling
stops as soon as a roll satisfies all of them. Every roll is appended to
`log.txt`, and the running best-sum roll is also written there for reference.

## Requirements

The target game server is expected to be listening on `127.0.0.1:25042`
before the roller is started.

## Development

This repo uses [mise](https://mise.jdx.dev) to pin the Go toolchain and
[just](https://github.com/casey/just) to run common tasks.

```sh
mise install     # install the pinned Go version
just build        # build ./roller
just run          # go run .
just check        # gofmt + go vet + go test
just clean         # remove build artifacts and log.txt
```

Run `just` with no arguments to list all available recipes.
