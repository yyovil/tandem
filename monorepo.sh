tui() {
  pushd ./tui && go run . ; popd
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