#!/bin/sh

set -eu

repository="drahil/ink"
install_dir="${INK_INSTALL_DIR:-/usr/local/bin}"

fail() {
    printf 'ink installer: %s\n' "$1" >&2
    exit 1
}

command -v curl >/dev/null 2>&1 || fail "curl is required"
command -v tar >/dev/null 2>&1 || fail "tar is required"

case "$(uname -s)" in
    Linux) os="linux" ;;
    Darwin) os="darwin" ;;
    *) fail "unsupported operating system: $(uname -s)" ;;
esac

case "$(uname -m)" in
    x86_64|amd64) arch="amd64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) fail "unsupported CPU architecture: $(uname -m)" ;;
esac

if [ -n "${INK_VERSION:-}" ]; then
    tag="$INK_VERSION"
else
    latest_url="$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${repository}/releases/latest")" || \
        fail "could not determine the latest release"
    tag="${latest_url##*/}"
fi

case "$tag" in
    v*) version="${tag#v}" ;;
    *) version="$tag"; tag="v${tag}" ;;
esac

case "$version" in
    *[!0-9A-Za-z._-]*|'') fail "invalid release version: $version" ;;
esac

archive="ink_${version}_${os}_${arch}.tar.gz"
download_root="${INK_DOWNLOAD_ROOT:-https://github.com/${repository}/releases/download/${tag}}"
temporary_dir="$(mktemp -d 2>/dev/null || mktemp -d -t ink-install)"
trap 'rm -rf "$temporary_dir"' EXIT HUP INT TERM

printf 'Downloading ink %s for %s/%s...\n' "$version" "$os" "$arch"
curl -fsSL "${download_root}/${archive}" -o "${temporary_dir}/${archive}" || \
    fail "could not download ${archive}"
curl -fsSL "${download_root}/checksums.txt" -o "${temporary_dir}/checksums.txt" || \
    fail "could not download checksums.txt"

expected_checksum="$(awk -v file="$archive" '$2 == file { print $1 }' "${temporary_dir}/checksums.txt")"
[ -n "$expected_checksum" ] || fail "release checksum for ${archive} was not found"

if command -v sha256sum >/dev/null 2>&1; then
    actual_checksum="$(sha256sum "${temporary_dir}/${archive}" | awk '{ print $1 }')"
elif command -v shasum >/dev/null 2>&1; then
    actual_checksum="$(shasum -a 256 "${temporary_dir}/${archive}" | awk '{ print $1 }')"
else
    fail "sha256sum or shasum is required to verify the download"
fi

[ "$actual_checksum" = "$expected_checksum" ] || fail "download checksum did not match"

tar -xzf "${temporary_dir}/${archive}" -C "$temporary_dir"
[ -f "${temporary_dir}/ink" ] || fail "the downloaded archive did not contain ink"

if [ -d "$install_dir" ] && [ -w "$install_dir" ]; then
    install -m 755 "${temporary_dir}/ink" "${install_dir}/ink"
else
    command -v sudo >/dev/null 2>&1 || \
        fail "cannot write to ${install_dir}; set INK_INSTALL_DIR to a writable directory"
    sudo mkdir -p "$install_dir"
    sudo install -m 755 "${temporary_dir}/ink" "${install_dir}/ink"
fi

printf 'Installed ink %s to %s/ink\n' "$version" "$install_dir"
printf 'Run "ink" from the directory you want to edit.\n'
