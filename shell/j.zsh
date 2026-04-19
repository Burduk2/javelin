function j() {
  local dir
  export EXE_ROOT="__EXE_ROOT__"
  dir="$("__BIN_PATH__" "$@")" || return 1
  cd "$dir" || return 1
}
