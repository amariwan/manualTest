#!/usr/bin/env bash

# arguments: ERROR_TEXT
usage()
{
    if [ -n "$1" ]; then
        echo "$(basename $0): $1" 1>&2
        echo
    fi

    echo "Usage: $0 [-v <VERSION> -b <BUILD>] [-s] [-h]"
    echo
    echo "Build application for web platform"
    echo
    echo "    Options:"
    echo "      -v <VERSION>   Application version in the form x.y.z"
    echo "      -b <BUILD>     Build number"
    echo "      -s             Create source maps"
    echo "      -H             Allow editing host field on login page"
}

SOURCE_MAPS=0
RW_HOST=0
VERSION=$(whoami)
BUILD=$(git branch --show-current 2>/dev/null)

while getopts "v:b:sHh" flag ; do
    case "${flag}" in
        v) VERSION=${OPTARG};;
        b) BUILD=${OPTARG};;
        s) SOURCE_MAPS=1;;
        H) RW_HOST=1;;
        h) usage; exit 0;;
        --) break;;
        *) echo "Invalid option ${flag}"; exit 1;;
    esac
done

set -x
set -e

if [ ${SOURCE_MAPS} -eq 1 ]; then
    SOURCE_MAPS_PARAM="--source-maps"
else
    rm -f build/web/main.dart.js.map
fi

flutter build web \
    --dart-define=FLUTTER_WEB_CANVASKIT_URL=canvaskit/ \
    --dart-define=build-date=$(date +%s) \
    --dart-define=build-name=${VERSION} \
    --dart-define=build-number=${BUILD} \
    --dart-define=rw-host=${RW_HOST}
    ${SOURCE_MAPS_PARAM} 

# the parameter '-i' sed on MacOS requires an additional argument
if [ "$(uname -s)" = Darwin ]; then DUMMY_OPT=".orig"; fi
sed \
-i ${DUMMY_OPT} build/web/main.dart.js

rm -f build/web/main.dart.js.orig
rm -f manuelle_test_web.zip.zip
rm -rf build/web/canvaskit/profiling

cd build/web
zip -r ../../manuelle_test_web.zip *
cd ../..
