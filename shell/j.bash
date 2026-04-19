j() {
  local script_dir exe_root dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  exe_root="$(cd "$script_dir/../bin" && pwd)"
  dir="$(EXE_ROOT="$exe_root" "$exe_root/javelin" "$@")" || return 1
  cd "$dir" || return 1
}
