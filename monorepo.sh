tandem() {
  pushd ./tui && go run . ; popd
}

# agentsup() {
#   pushd ./agents && ag ws up ; popd
# }

agentsup() {
  source agents/.env && fastapi dev --port 8000 agents/api/main.py
}
agentsdown() {
  pushd ./agents && ag ws down ; popd
}

logs() {
  cat tui/debug.log
}
