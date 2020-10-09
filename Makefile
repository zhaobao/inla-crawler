.PHONY: run-comic run-novel run-med run-cozy run-anime \
	run-readnovelfull run-gds run-wiz run-comic-water \
	run-netease-us run-translate \
	run-netease-ng

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

run-wiz:
	go run tasks/video/wizlimited/main.go

run-comic-water:
	go run factory/watermark/comic/main.go

run-netease-us:
	go run tasks/music/netease/us/main.go

run-netease-ng:
	go run tasks/music/netease/ng/main.go

run-translate:
	go run tasks/translation/game/main.go