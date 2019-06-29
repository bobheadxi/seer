
.PHONY: typescript
typescript:
	tscriptify -package=go.bobheadxi.dev/seer/riot -target=web/src/api/types.ts riot/api.go

.PHONY: lint
lint:
	./.scripts/lint.sh