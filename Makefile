# If the first argument is "terraform"...
ifeq (terraform,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "terraform"
  TF_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(TF_ARGS):;@:)
endif

build: 
	$(MAKE) -C storage build
	$(MAKE) -C common build
	$(MAKE) -C lambdas/fetch-changed build
	$(MAKE) -C lambdas/fetch-schema.org  build
	$(MAKE) -C lambdas/merge-schema.org  build
	$(MAKE) -C lambdas/transform-parsoid build

clean:
	$(MAKE) -C storage clean
	$(MAKE) -C common clean
	$(MAKE) -C lambdas/fetch-changed clean
	$(MAKE) -C lambdas/fetch-schema.org clean
	$(MAKE) -C lambdas/merge-schema.org clean
	$(MAKE) -C lambdas/transform-parsoid clean

deploy: 
	$(MAKE) -C lambdas/fetch-changed deploy
	$(MAKE) -C lambdas/fetch-schema.org  deploy
	$(MAKE) -C lambdas/merge-schema.org  deploy
	$(MAKE) -C lambdas/transform-parsoid deploy

test: 
	$(MAKE) -C storage test
	$(MAKE) -C common test
	$(MAKE) -C lambdas/fetch-changed test
	$(MAKE) -C lambdas/fetch-schema.org  test
	$(MAKE) -C lambdas/merge-schema.org  test
	$(MAKE) -C lambdas/transform-parsoid test

terraform:
	@terraform -chdir=iac $(TF_ARGS)
	
tflint:
	tflint -c ./iac/.tflint.hcl ./iac

.PHONY: build clean deploy test terraform tflint
