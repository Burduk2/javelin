function j() {
  local script_dir exe_root dir
  script_dir="${${(%):-%N}:A:h}"
  exe_root="${script_dir:h}/bin"
  dir="$(EXE_ROOT="$exe_root" "$exe_root/javelin" "$@")" || return 1
  cd "$dir" || return 1
}
