HEROKU_APP=seer-engine

all: docs
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

.PHONY: dev-web
dev-web:
	cd web ; npm run serve

.PHONY: docs
docs:
	cd web ; npm run build

.PHONY: server
server:
	docker-compose up --build

# https://developer.riotgames.com/
.PHONY: riot-token
riot-token:
	heroku config:set -a $(HEROKU_APP) RIOT_API_TOKEN=$(RIOT_API_TOKEN)

.PHONY: check-server
check-server:
	@heroku apps:info --app $(HEROKU_APP)
	@echo '=== server status'
	@curl https://$(HEROKU_APP).herokuapp.com/status

.PHONY: prod-redis
prod-redis:
	heroku config -a $(HEROKU_APP) | grep REDIS_URL

.PHONY: prod-jobs
prod-jobs:
	go run . -mode=jobs-ui-only
