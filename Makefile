SHELL:=/bin/zsh

check-modd-exists:
	@modd --version > /dev/null

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf

