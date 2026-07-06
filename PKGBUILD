# Maintainer: Foxemsx
pkgname=speed
pkgver=1.0.0
pkgrel=1
pkgdesc="Internet speed test in your terminal — download & upload"
arch=('x86_64' 'aarch64')
url="https://github.com/Foxemsx/speed"
license=('MIT')
depends=()
makedepends=('go')
source=("$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$pkgname-$pkgver"
  go build -trimpath -ldflags "-s -w" -o "$pkgname" .
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
