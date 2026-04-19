function j
    set -l script_dir (path dirname (status filename))
    set -l exe_root (path resolve "$script_dir/../bin")
    set -l dir (env EXE_ROOT="$exe_root" "$exe_root/javelin" $argv)
    or return 1
    cd "$dir"
end
