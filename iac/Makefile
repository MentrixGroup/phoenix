AWS_DEFAULT_REGION ?= us-east-1

documentation:
	@terraform-docs markdown . > README.md

validate:
	@export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION}"; \
	$(MAKE) init; \
	terraform validate

init:
	@terraform init -get -upgrade