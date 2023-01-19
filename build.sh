#!/usr/bin/env bash

PROG_NAME=$(basename "$0")

function process_command() {
  while [[ $# -gt 0 ]]; do
    case $1 in
    -h | --help | help)
      show_command_help
      return 0
      ;;
    -d | --dry-run)
      DRY_RUN="YES"
      shift # past command argument
      ;;
    *)
      echo "ERROR: Unknown command argument $1"
      show_command_help
      return 1
      ;;
    esac
  done

  build
  return 0
}

function show_command_help() {
  local HELP
  read -r -d '' HELP <<-EOF
		Description:
		  Used to build the CI binaries for Linux, Windows and macOs

		Usage:
		  ${PROG_NAME} [-d|--dry-run]
		  ${PROG_NAME} [-h|--help]

		Flags:
		  -d, --dry-run    Compares temp to committed binaries, but does not replace them
		  -h, --help       help for ${PROG_NAME}

		Temp binaries are available in the ./tmp directory
EOF
  echo "${HELP}"
}

function build() {
  local go_oses=(
    "GOOS=darwin GOARCH=amd64"
    "GOOS=darwin GOARCH=arm64"
    "GOOS=linux GOARCH=amd64"
    "GOOS=linux GOARCH=arm64"
    "GOOS=windows GOARCH=amd64"
    "GOOS=windows GOARCH=arm64"
  )
  local binary_suffixes=(
    "darwin-amd64"
    "darwin-arm64"
    "linux-amd64"
    "linux-arm64"
    "windows-amd64.exe"
    "windows-arm64.exe"
  )

  if [[ ${#go_oses[@]} -ne ${#binary_suffixes[@]} ]]; then
    echo "ERROR: List of oses and binary suffixes are mismatched!"
    exit 1
  fi

  printf "🚧 Starting!%s" "$([[ -n "${DRY_RUN}" ]] && echo " -- DRY RUN!!! (No binaries will be updated)")"
  echo
  echo
  rm -rf ./tmp

  for i in "${!go_oses[@]}"; do
    local binary_suffix="${binary_suffixes[$i]}"
    local new_binary_path="tmp/get-app-token-${binary_suffix}"
    local old_binary_path="bin/get-app-token-${binary_suffix}"
    local build_cmd="${go_oses[$i]} go build -ldflags='-s -w' -o $new_binary_path"

    echo "Building binary for ${binary_suffix}"
    eval "${build_cmd}"

    if [[ -f ${old_binary_path} ]]; then
      local new_binary_sha
      local old_binary_sha

      new_binary_sha="$(get_sha "${new_binary_path}")"
      old_binary_sha="$(get_sha "${old_binary_path}")"

      if [[ "${new_binary_sha}" == "${old_binary_sha}" ]]; then
        echo "  Binary for ${binary_suffix} is up to date. 👍"
      else
        echo "  Replacing binary for ${binary_suffix} with newer version"
        echo "  🔻 (old) ${old_binary_sha}"
        echo "  ✅ (new) ${new_binary_sha}"
        move_file "${new_binary_path}" "${old_binary_path}"
      fi
    else
      echo " Adding binary for ${binary_suffix}"
      move_file "${new_binary_path}" "${old_binary_path}"
    fi

    echo "✓ Build complete for ${binary_suffixes[$i]}: ${old_binary_path}"
  done

  echo
  printf "🏁 Done!%s" "$([[ -n ${DRY_RUN} ]] && echo " -- DRY RUN!!! (No binaries were updated)")"
  echo
}

function get_sha() {
  local binary_path="$1"
  local sha_str
  local sha

  sha_str="$(sha256sum --binary --tag "${binary_path}")"
  sha="$(echo "${sha_str}" | sed -n -e 's/^.* = //p')"

  echo "$sha"
}

function move_file() {
  local new_binary_path="${1}"
  local old_binary_path="${2}"

  if [[ -z "${DRY_RUN}" ]]; then
    [[ ! -d  "./bin" ]] && mkdir "./bin"

    mv -f "${new_binary_path}" "${old_binary_path}"
  fi
}

process_command "$@"
