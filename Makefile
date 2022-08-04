push:
	cmd/git_push.sh

composeUp:
	docker-compose  up --build

composeDown:
	docker-compose  down --remove-orphans

.PHONY: push composeUp composeDown

