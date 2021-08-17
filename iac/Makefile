MODULES := $(wildcard */.)
AWS_DEFAULT_REGION ?= us-east-1

documentation:
	@terraform-docs markdown . > README.md
	@for dir in $(MODULES); do \
		printf "\n\nGenerating documentation for $$dir\n"; \
		$(MAKE) -C $$dir documentation; \
	done

validate:
	@export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}"; \
	$(MAKE) init; \
	terraform validate

init:
	@terraform init -get -upgrade