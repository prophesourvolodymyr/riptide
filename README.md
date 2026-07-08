# speed

A polished terminal internet speed test written in Go.

It measures **download**, **upload**, and **latency/ping** with a centered
card UI, live progress feedback, and compact speed history graphs.

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

## Download

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

## License

See the repository license file.
