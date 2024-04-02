#!/usr/bin/env sh

set -eu

push() {
    if ! git remote get-url "${1}" 1> /dev/null 2> /dev/null; then
        git remote add "${1}" "${2}"
    fi
    git push "${1}" main
}

cd "$(dirname "$(realpath "$0")")/.."

push origin 'git@git.sr.ht:~hub/semver'
push github 'git@github.com:nothub/semver.git'
