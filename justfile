set quiet

run *args:
  -go run ./src {{args}}

build:
  go build -o bin/javelin ./src

install shell="auto":
  #!/usr/bin/env bash
  set -euo pipefail
  repo_root="$(pwd)"
  bin_dir="$repo_root/bin"
  shell_name="{{shell}}"
  if [ "$shell_name" = "auto" ]; then
    shell_name="$(basename "${SHELL:-}")"
  fi
  go build -o "$bin_dir/javelin" ./src
  case "$shell_name" in
    zsh)
      wrapper="$repo_root/shell/j.zsh"
      rc_file="$HOME/.zshrc"
      ;;
    bash)
      wrapper="$repo_root/shell/j.bash"
      rc_file="$HOME/.bashrc"
      ;;
    fish)
      wrapper="$repo_root/shell/j.fish"
      rc_file="$HOME/.config/fish/config.fish"
      mkdir -p "$(dirname "$rc_file")"
      ;;
    *)
      echo "unsupported shell: $shell_name"
      echo "supported shells: zsh, bash, fish"
      exit 1
      ;;
  esac
  touch "$rc_file"
  source_line="source $wrapper"
  if ! grep -Fqx "$source_line" "$rc_file"; then
    printf '\n%s\n' "$source_line" >> "$rc_file"
    echo "added j wrapper to $rc_file"
  else
    echo "j wrapper already installed in $rc_file"
  fi
  echo "built $bin_dir/javelin"
