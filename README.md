# speed

Internet speed test in your terminal — download & upload. Powered by [fast.com](https://fast.com).

![demo](demo.gif)

## Install

### One-liner (any Linux/macOS)

```sh
curl -sSL https://raw.githubusercontent.com/Foxemsx/speed/main/install.sh | bash
```

### Go

```sh
go install github.com/Foxemsx/speed@latest
```

### Homebrew (macOS / Linux)

```sh
brew install Foxemsx/speed/speed
```

### AUR (Arch Linux)

```sh
yay -S speed
```

## Usage

```sh
speed
```

Run a download test followed by an upload test. Results are displayed with a live sparkline graph and peak speeds.

## How it works

`speed` measures your connection against the nearest Netflix Open Connect servers (via fast.com). It runs parallel downloads and uploads for 10 seconds each, reporting speed in Mbps (auto-scaling to Gbps for fast connections).

## License

[MIT](LICENSE)
