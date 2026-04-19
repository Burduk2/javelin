function j
    set -gx EXE_ROOT "__EXE_ROOT__"
    set -l dir ("__BIN_PATH__" $argv)
    or return 1
    cd "$dir"
end
