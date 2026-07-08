# speed

A polished terminal internet speed test written in Go, inspired by
[maaslalani/fast](https://github.com/maaslalani/fast) but with **upload
testing** and a significantly more refined visual design.

It measures **download**, **upload**, and **latency/ping** using the same
approach as fast.com: it fetches a list of nearby Netflix Open Connect CDN
targets, opens several parallel HTTP connections to saturate the link, counts
bytes transferred over a fixed duration, then probes latency with a final
round-trip.

![card](https://img.shields.io/badge/terminal-speed%20test-⚡-39d0d8)

## Features

- **Download + upload + ping** in one run.
- **Parallel connections** (~5) to saturate your link, like fast.com.
- A **centered, rounded-border card** that reflows on resize and looks good at
  any terminal size.
- **Custom full-screen background** via `--bg <hex>` (falls back to the
  terminal's default background when omitted).
- **Distinct accent colors** for download (teal) vs upload (amber), with a
  small reskinnable theme.
- **Compact live sparklines** beside each speed readout, updated every
  ~100 ms.
- **Smooth interpolation** (lerp) so the displayed numbers and bars glide
  instead of snapping.
- **Phase progression**: finding servers (spinner) → download → upload →
  latency → one-line summary with peak values.
- **Graceful errors**: no internet / network failures show a clear message,
  not a stack trace.
- Quits cleanly on **q / esc / ctrl+c**, cancelling in-flight transfers.

## Install

```sh
go install github.com/rohli/speed@latest
```

Or build from a clone:

```sh
git clone https://github.com/rohli/speed
cd speed
go build -o speed .
```

Requires Go 1.23+.

## Usage

```sh
speed                 # run the test with the terminal's default background
speed --bg "#0d1117"  # run with a custom full-screen background color
speed --bg 0d1117     # leading '#' is optional
```

### Flags

| Flag     | Description                                              | Default |
| -------- | -------------------------------------------------------- | ------- |
| `--bg`   | Full-screen background color as hex (`#rrggbb`). Omit to use the terminal's native background. | _(none)_ |
| `--theme`| Color theme name. Currently only `default`. Reserved for future palettes. | `default` |

After the run completes, press **q**, **esc**, or **ctrl+c** to quit (you can
quit at any time to abort an in-progress test).

## Architecture

The code is split into small, dependency-light files:

| File           | Responsibility |
| -------------- | -------------- |
| `main.go`      | CLI entry point: flag parsing (`--bg`), background-color validation, and `tea.NewProgram` setup with the alternate screen. |
| `speed.go`     | Measurement engine. Discovers targets from the fast.com config API, runs parallel download (`GET`) and upload (`POST`) phases with atomic byte counters, probes latency, and streams `Phase`/`Sample`/`Result` over channels. |
| `model.go`     | The bubbletea `Model`: event loop (`Update`), the `View` (centered bordered card), the 100 ms `tick` refresh, **lerp-based animation** of displayed values toward measured targets, and the timed-phase watchdog/timer UI. |
| `sparkline.go` | A compact single-row sparkline renderer for speed-over-time. |
| `theme.go`     | The reskinnable `Theme` struct (bg/fg/accent colors) and the default dark palette. |

### Data flow

1. `main` starts bubbletea, which `Init`s the spinner, the tick loop, and the
   background `Run` goroutine.
2. `Run` fetches target URLs, then runs download → upload → latency, emitting
  `Phase` transitions and stamped `Sample`s (instantaneous rates in bytes/sec)
  on channels the model drains each tick.
3. Each tick the model **lerps** the displayed download/upload numbers toward
  the latest measured rate, advances the compact sparkline, and updates the
  live seconds counter for the active timed phase.
4. The UI also keeps a local watchdog so it can transition between phases even
  if the background runner stalls or misses an event.
5. On completion a `Result` (with Mbps + peak values) arrives and the card
  switches to the final summary phase.

### Tradeoffs

- **Animation smoothing — lerp, not a spring.** I used a simple per-tick
  linear interpolation (`lerp(current, target, 0.18)`) rather than a
  harmonica spring. Lerp is trivially correct, needs no velocity state, and
  looks fluid at a 100 ms tick. A spring would add overshoot/bounce, which
  reads as "playful" but can feel less precise for a measurement tool. The
  `animFactor` constant makes it easy to swap in a spring later.
- **Timed-phase fallback.** The download and upload phases use a local timer
  in the UI as a fallback, so the card cannot sit forever on a stale phase
  label if a background event is delayed.
- **Default colors.** Picked a deep-slate background (`#0d1117`) with a cool
  teal for *incoming* (download) and a warm amber for *outgoing* (upload) —
  the warm/cool split makes the two directions instantly distinguishable, and
  muted greys carry structure. All text uses `lipgloss.AdaptiveColor` so it
  stays readable on both light and dark `--bg` values.
- **Single token.** The fast.com v2 config API requires a public app token
  (the one shipped in fast.com's web bundle). It is not secret, but Netflix
  could rotate or block it; if the server list fails to load you'll get a
  clear "could not reach speed-test servers" message rather than a crash.
- **Upload method.** fast.com uploads random `application/octet-stream` POST
  bodies to the same CDN targets. I reuse a single 256 KB random buffer posted
  repeatedly for the duration, which approximates fast.com's increasing-size
  uploads well enough for a steady-state speed reading.
- **Latency** is a single round-trip probe to one target after transfers
  finish (not an averaged min-of-N), keeping the run short and simple.
- **Harmonica** (`charmbracelet/harmonica`) is listed as an option in the
  brief but intentionally not used, for the reason above — keeps the binary
  dependency-light and the motion predictable.
