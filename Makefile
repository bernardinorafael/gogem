MODULES = apiutil cache crypto dbutil fault function httputil logger pagination queue server uid

## help: show available commands
.PHONY: help
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'

## vet: run go vet on all modules
.PHONY: vet
vet:
	@for mod in $(MODULES); do \
		echo "vet pkg/$$mod"; \
		cd pkg/$$mod && go vet ./... && cd ../..; \
	done

## test: run tests on all modules
.PHONY: test
test:
	@for mod in $(MODULES); do \
		echo "test pkg/$$mod"; \
		cd pkg/$$mod && go test ./... && cd ../..; \
	done

## tidy: run go mod tidy on all modules
.PHONY: tidy
tidy:
	@for mod in $(MODULES); do \
		cd pkg/$$mod && go mod tidy && cd ../..; \
	done

## fmt: format all Go files
.PHONY: fmt
fmt:
	gofmt -w .

## tags: list the latest tag for each module
.PHONY: tags
tags:
	@for mod in $(MODULES); do \
		tag=$$(git tag -l "pkg/$$mod/v*" --sort=-v:refname | head -1); \
		if [ -z "$$tag" ]; then \
			printf "  %-12s  (no tags)\n" "$$mod"; \
		else \
			printf "  %-12s  %s\n" "$$mod" "$$tag"; \
		fi; \
	done

## changed: show modules changed since their last tag
.PHONY: changed
changed:
	@found=false; \
	for mod in $(MODULES); do \
		tag=$$(git tag -l "pkg/$$mod/v*" --sort=-v:refname | head -1); \
		if [ -z "$$tag" ]; then \
			printf "  %-12s  (never tagged)\n" "$$mod"; \
			found=true; \
		else \
			changes=$$(git diff --name-only "$$tag"..HEAD -- "pkg/$$mod/"); \
			if [ -n "$$changes" ]; then \
				count=$$(echo "$$changes" | wc -l | tr -d ' '); \
				printf "  %-12s  %s files changed since %s\n" "$$mod" "$$count" "$$tag"; \
				found=true; \
			fi; \
		fi; \
	done; \
	if [ "$$found" = false ]; then \
		echo "  No modules changed since last tags."; \
	fi

## tag: create a version tag for a module. Usage: make tag m=fault b=patch
.PHONY: tag
tag:
ifndef m
	$(error Usage: make tag m=<module> b=<patch|minor|major>)
endif
ifndef b
	$(error Usage: make tag m=<module> b=<patch|minor|major>)
endif
	@if ! echo "$(MODULES)" | grep -qw "$(m)"; then \
		echo "Error: unknown module '$(m)'"; \
		echo "Available: $(MODULES)"; \
		exit 1; \
	fi
	@if ! echo "patch minor major" | grep -qw "$(b)"; then \
		echo "Error: bump must be patch, minor, or major"; \
		exit 1; \
	fi
	@current_tag=$$(git tag -l "pkg/$(m)/v*" --sort=-v:refname | head -1); \
	if [ -z "$$current_tag" ]; then \
		new_version="0.1.0"; \
	else \
		version=$${current_tag#pkg/$(m)/v}; \
		major=$$(echo $$version | cut -d. -f1); \
		minor=$$(echo $$version | cut -d. -f2); \
		patch=$$(echo $$version | cut -d. -f3); \
		case "$(b)" in \
			major) new_version="$$((major + 1)).0.0" ;; \
			minor) new_version="$$major.$$((minor + 1)).0" ;; \
			patch) new_version="$$major.$$minor.$$((patch + 1))" ;; \
		esac; \
	fi; \
	new_tag="pkg/$(m)/v$$new_version"; \
	if [ -n "$$current_tag" ]; then \
		echo "  $$current_tag -> $$new_tag"; \
	else \
		echo "  (new) -> $$new_tag"; \
	fi; \
	echo ""; \
	printf "  Create tag $$new_tag? [y/N] "; \
	read confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		git tag "$$new_tag"; \
		echo "  Tag $$new_tag created."; \
		echo ""; \
		printf "  Push $$new_tag to origin? [y/N] "; \
		read push_confirm; \
		if [ "$$push_confirm" = "y" ] || [ "$$push_confirm" = "Y" ]; then \
			git push origin "$$new_tag"; \
			echo "  Pushed $$new_tag."; \
		else \
			echo "  Tag created locally. Push later with: git push origin $$new_tag"; \
		fi; \
	else \
		echo "  Aborted."; \
	fi
