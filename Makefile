
all:
	rm -rf seer
	go build

.PHONY: typescript
typescript:
	tscriptify -package=go.bobheadxi.dev/seer/riot -target=web/src/api/types.ts riot/api.go

.PHONY: lint
lint:
	./.scripts/lint.sh

.PHONY: web
web:
	cd web ; npm run serve:prod

.PHONY: server
server:
	docker-compose up
