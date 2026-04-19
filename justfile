set quiet

run *args:
  -go run ./src {{args}}

build:
  GOOS=darwin GOARCH=arm64 go build -o ./bin/javelin-darwin-arm64 ./src
  GOOS=darwin GOARCH=amd64 go build -o ./bin/javelin-darwin-amd64 ./src

  GOOS=linux GOARCH=amd64 go build -o ./bin/javelin-linux-amd64 ./src
  GOOS=linux GOARCH=arm64 go build -o ./bin/javelin-linux-arm64 ./src

  GOOS=windows GOARCH=amd64 go build -o ./bin/javelin-windows-amd64.exe ./src
  GOOS=windows GOARCH=arm64 go build -o ./bin/javelin-windows-arm64.exe ./src

install shell="auto":
  #!/usr/bin/env bash
  set -euo pipefail
  repo_root="$(pwd)"
  shell_name="{{shell}}"
  install_root="${XDG_DATA_HOME:-$HOME/.local/share}/javelin"
  bin_dir="$install_root/bin"
  config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/javelin"
  wrapper_dest=""
  template=""
  rc_file=""
  installed_bin="$bin_dir/javelin"
  escape_sed_replacement() {
    printf '%s' "$1" | sed -e 's/[\/&]/\\&/g'
  }
  if [ "$shell_name" = "auto" ]; then
    shell_name="$(basename "${SHELL:-}")"
  fi
  mkdir -p "$bin_dir" "$config_dir"
  go build -o "$installed_bin" ./src
  case "$shell_name" in
    zsh)
      template="$repo_root/shell/j.zsh"
      wrapper_dest="$config_dir/j.zsh"
      rc_file="$HOME/.zshrc"
      ;;
    bash)
      template="$repo_root/shell/j.bash"
      wrapper_dest="$config_dir/j.bash"
      rc_file="$HOME/.bashrc"
      ;;
    fish)
      template="$repo_root/shell/j.fish"
      wrapper_dest="$config_dir/j.fish"
      rc_file="$HOME/.config/fish/config.fish"
      mkdir -p "$(dirname "$rc_file")"
      ;;
    *)
      echo "unsupported shell: $shell_name"
      echo "supported shells: zsh, bash, fish"
      exit 1
      ;;
  esac
  exe_root_escaped="$(escape_sed_replacement "$bin_dir")"
  bin_path_escaped="$(escape_sed_replacement "$installed_bin")"
  sed \
    -e "s/__EXE_ROOT__/$exe_root_escaped/g" \
    -e "s/__BIN_PATH__/$bin_path_escaped/g" \
    "$template" > "$wrapper_dest"
  touch "$rc_file"
  source_line="source $wrapper_dest"
  if ! grep -Fqx "$source_line" "$rc_file"; then
    printf '\n%s\n' "$source_line" >> "$rc_file"
    echo "added j wrapper to $rc_file"
  else
    echo "j wrapper already installed in $rc_file"
  fi
  echo "installed $installed_bin"
  echo "rendered $wrapper_dest"
