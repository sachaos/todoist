workflow "New workflow" {
  on = "push"
  resolves = ["golang"]
}

action "golang" {
  uses = "golang"
}
