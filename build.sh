#!/usr/bin/env bash

package=$1

ARTIFACTS_DIR=artifacts
VERSION=$CI_COMMIT_TAG
BUILDTIME=$(date +"%Y.%m.%d.%H%M%S")
REVISION=$(git log --pretty=format:"%h" -n 1)
BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [ -z "$VERSION" ]; then
    VERSION=$BRANCH'-'$REVISION
fi

if [[ -z "$package" ]]; then
  echo "Usage: $0 <package-name>"
  exit 1
fi

package_split=(${package//\// })
package_name=${package_split[-1]}
echo "Package name: ${package_name}"

platforms=("windows/amd64" "linux/amd64" "windows/386")
for platform in "${platforms[@]}"; do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name=$package_name
  archive_name=$package_name'-'$GOOS'-'$GOARCH'-'$VERSION
  if [ $GOOS = "windows" ]; then
    output_name+='.exe'
    archive_name+='.zip'
  else 
    archive_name+='.tgz'
  fi

  # env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
  env GOOS=$GOOS GOARCH=$GOARCH go build \
  -ldflags "-extldflags '-static' \
  -X ozz-ms/ozz-srv/cmd.VERSION=$VERSION \
  -X ozz-ms/ozz-srv/cmd.BUILT=$BUILDTIME \
  -X ozz-ms/ozz-srv/cmd.REVISION=$REVISION \
  -X ozz-ms/ozz-srv/cmd.BRANCH=$BRANCH \
  " \
  -o ./$ARTIFACTS_DIR/$output_name $package
  if [ $? -ne 0 ]; then
    echo 'Error occured'
    exit 1
  fi
  echo "Output name: ${output_name}"

  if [ $GOOS = "windows" ]; then
    zip -m -j ./$ARTIFACTS_DIR/$archive_name ./$ARTIFACTS_DIR/$output_name
  else
    tar zcf ./$ARTIFACTS_DIR/$archive_name -C ./$ARTIFACTS_DIR $output_name
    rm ./$ARTIFACTS_DIR/$output_name
  fi
done
