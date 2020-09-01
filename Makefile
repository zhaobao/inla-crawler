.PHONY: run-comic run-novel run-med run-cozy run-anime run-readnovelfull run-gds

run-comic:
	go run tasks/qgxymdmz/main.go

run-novel:
	go run tasks/novel/qgxymdmz/main.go

run-med:
	go run tasks/meditation/tide/main.go

run-cozy:
	go run tasks/music/cozy/main.go

run-anime:
	go run tasks/anime/popsanime/main.go

run-readnovelfull:
	go run tasks/novel/readnovelfull/main.go

run-gds:
	go run tasks/novel/gdsbook/main.go