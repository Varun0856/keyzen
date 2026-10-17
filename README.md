# keyzen

A fast, minimal typing test for your terminal, built in Go with bubbletea.
Type a random set of words, see correct and wrong characters highlighted live,
and get your WPM and accuracy when you finish. Timed modes, saved best scores,
and a single binary you can install with `go install`.

## Install

```bash
go install github.com/Varun0856/keyzen/cmd/keyzen@latest
```

Or build from source:

```bash
git clone https://github.com/Varun0856/keyzen
cd keyzen
go build -o keyzen ./cmd/keyzen
```

## Usage

```bash
keyzen          # 30 second test
keyzen -t 15    # 15 second test
keyzen -t 60    # 60 second test
```

| Key       | Action          |
| --------- | --------------- |
| any key   | start the timer |
| backspace | fix a mistake   |
| tab       | restart         |
| esc       | quit            |

## How it's measured

- **WPM** counts only correct characters, with 5 characters = 1 word.
- **Accuracy** counts every keystroke, so correcting a mistake with
  backspace still counts against you.
- Best scores per mode are saved to `~/.config/keyzen/scores.json`
  (the location follows your OS's config directory).

## Project layout

```
cmd/keyzen/        entry point and bubbletea UI
internal/typing/   session state, WPM and accuracy (no UI code, fully tested)
internal/words/    embedded word list and random generation
internal/storage/  best-score persistence
```

## License

MIT
