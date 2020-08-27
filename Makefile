.PHONY: run-comic run-novel run-med run-cozy

run-comic:
	go run tasks/qgxymdmz/main.go

run-novel:
	go run tasks/novel/qgxymdmz/main.go

run-med:
	go run tasks/meditation/tide/main.go

run-cozy:
	go run tasks/music/cozy/main.go