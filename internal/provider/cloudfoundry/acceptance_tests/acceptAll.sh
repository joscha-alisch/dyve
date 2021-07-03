for f in ./*.actual.json; do mv "$f" "$(echo "$f" | sed s/actual/accepted/)"; done
