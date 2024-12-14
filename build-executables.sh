VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

platforms=(
"darwin/amd64"
"darwin/arm64"
"linux/amd64"
"linux/arm"
"linux/arm64"
"windows/amd64"
)

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}

  os=$GOOS
  if [ $os = "darwin" ]; then
    os="macOS"
  fi

  output_name="dev-setup-${version}-${os}-${GOARCH}"
  if [ $os = "windows" ]; then
    output_name+='.exe'
  fi

  echo "Building release/$output_name..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X github.com/DragonSecurity/dev-setup/main.Version=${VERSION} -X github.com/DragonSecurity/dev-setup/main.Commit=${COMMIT} -X github.com/DragonSecurity/dev-setup/main.Date=${DATE}" -o release/$output_name
  if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting.'
    exit 1
  fi

  zip_name="dev-setup-${version}-${os}-${GOARCH}"
  pushd release > /dev/null
  if [ $os = "windows" ]; then
    zip $zip_name.zip $output_name
    rm $output_name
  else
    chmod a+x $output_name
    gzip $output_name
  fi
  popd > /dev/null
done