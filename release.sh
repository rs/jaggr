#!/bin/bash

set -e

VERSION=$TRAVIS_TAG
USER=rs
NAME=jaggr
DESC="JSON Aggregation CLI"

rm -rf dist
mkdir dist

envs=$(go tool dist list | grep -v android | grep -v darwin/arm | grep -v s390x | grep -v plan9/arm)

for env in $envs; do
    eval $(echo $env | tr '/' ' ' | xargs printf 'export GOOS=%s; export GOARCH=%s\n')

    GOOS=${env%/*}
    GOARCH=${env#*/}

    bin=$NAME
    if [ $GOOS == "windows" ]; then
        bin="$NAME.exe"
    fi

    mkdir -p dist

    echo "Building for GOOS=$GOOS GOARCH=$GOARCH"

    CGO_ENABLED=0 go build -o dist/$bin
    file=${NAME}_${VERSION}_${GOOS}_${GOARCH}.zip
    zip -q dist/$file -j dist/$bin
    rm -f dist/$bin
done

url=https://github.com/${USER}/${NAME}/archive/${VERSION}.tar.gz
darwin_amd64=${NAME}_${VERSION}_darwin_amd64.zip
darwin_386=${NAME}_${VERSION}_darwin_386.zip

cat << EOF > dist/homebrew.rb
class $(echo ${NAME:0:1} | tr '[a-z]' '[A-Z]')${NAME:1} < Formula
  desc "$DESC"
  homepage "https://github.com/${USER}/${NAME}"
  url "$url"
  sha256 "$(curl -s $url | shasum -a 256 | awk '{print $1}')"
  head "https://github.com/${USER}/${NAME}.git"

  if Hardware::CPU.is_64_bit?
    url "https://github.com/${USER}/${NAME}/releases/download/${VERSION}/${darwin_amd64}"
    sha256 "$(shasum -a 256 dist/${darwin_amd64} | awk '{print $1}')"
  else
    url "https://github.com/${USER}/${NAME}/releases/download/${VERSION}/${darwin_386}"
    sha256 "$(shasum -a 256 dist/${darwin_386} | awk '{print $1}')"
  end

  depends_on "go" => :build

  def install
    bin.install "$NAME"
  end
end
EOF
