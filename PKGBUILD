# This is an example PKGBUILD file. Use this as a start to creating your own,
# and remove these comments. For more information, see 'man PKGBUILD'.
# NOTE: Please fill out the license field for your package! If it is unknown,
# then please put 'unknown'.

# The following guidelines are specific to BZR, GIT, HG and SVN packages.
# Other VCS sources are not natively supported by makepkg yet.

# Maintainer: Souranil Chakraborty <chakrabortysoura21@gmail.com>
pkgname=NvFile 
pkgver=1.2.1
pkgrel=1
pkgdesc="A simple tui file explorer that can work with various text editors for all your developer needs."
arch=(x86_64)
url="https://github.com/Chakrabortysoura/NvFile"
license=('BSD-3')
groups=("utilities")
depends=()
makedepends=('go') # 'bzr', 'git', 'mercurial' or 'subversion'
provides=("${pkgname}")
conflicts=()
replaces=()
backup=()
options=()
install=
source=('https://github.com/Chakrabortysoura/NvFile/archive/refs/tags/v1.2.1.tar.gz')
noextract=()
sha256sums=('SKIP')
validpgpkeys=('SKIP')

# Please refer to the 'USING VCS SOURCES' section of the PKGBUILD man page for
# a description of each element in the source array.

prepare() {
	cd "$srcdir"
	rm v1.2.1.tar.gz
	mv NvFile-1.2.1 NvFile
	echo "Preparing packge"
}

build() {
	cd "$srcdir/NvFile"
	echo "Downloading go dependencies"
	go mod download 
	echo "Compiling the source binary"
	go build -o nvfile
}

check() {
	echo "skipping integrity checks"
}

package() {
	cd "$srcdir/NvFile"
	install -D nvfile "$pkgdir/usr/bin/nvfile"
}
