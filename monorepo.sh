tui() {
  pushd ./tui && go run . ; popd
}

run() {
  if declare -f "$1" > /dev/null; then
    "$1"
  else
    echo "Error: Unknown function '$1'"
    return 1
  fi
}

agentsup() {
  pushd ./agents && ag ws up ; popd
}

agentsdown() {
  pushd ./agents && ag ws down ; popd
}

logs() {
  cat tui/debug.log
}