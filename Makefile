SHELL:=/bin/zsh

migrate_up=go run main.go migrate --direction=up --step=0
migrate_down=go run main.go migrate --direction=down --step=0

check-modd-exists:
	@modd --version > /dev/null

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf


migrate:
	@if [ "$(DIRECTION)" = "" ] || [ "$(STEP)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run main.go migrate --direction=$(DIRECTION) --step=$(STEP);\
    fi

seed-user:
	@go run main.go seed-user