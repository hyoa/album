.PHONY: test unit-test mutation_test
test: cs phpstan unit-test ui-test

unit-test: vendor ## Run unit tests
	php vendor/bin/phpunit --testdox --group=unit --configuration phpunit.xml.dist

ui-test: vendor ## Run unit tests
	php vendor/bin/phpunit --testdox --group=ui --configuration phpunit.xml.dist

.PHONY: phpstan
phpstan: ## Check static analysis
	php -d memory_limit=-1 vendor/bin/phpstan analyse src tests --level=max

.PHONY: cs cs-fix cs-ci
cs: ## Check code style
	vendor/bin/php-cs-fixer fix --dry-run --stop-on-violation --diff

cs-fix: ## Fix code style
	vendor/bin/php-cs-fixer fix --allow-risky=yes

cs-ci: ## Run Continuous Integration code style check
	vendor/bin/php-cs-fixer fix --dry-run --using-cache=no --verbose --allow-risky=yes

.PHONY: deploy
deploy-dev: ##Deploy
	./build.sh && cd lambda_build && ./clean_dependencies.sh && serverless deploy --stage dev && cd ..

deploy-prod: ##Deploy
	./build.sh && cd lambda_build && serverless deploy --stage prod && cd ..
