.PHONY: run-comic run-novel run-med

run-comic:
	go run tasks/qgxymdmz/main.go

run-novel:
	go run tasks/novel/qgxymdmz/main.go

run-med:
	go run tasks/meditation/tide/main.go