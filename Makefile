.PHONY: tidy-all fmt-all vet-all build-all update-all clean-all

tidy-all:
	@echo "🔍 Tidying all Go modules..."
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "📦 go mod tidy in $$dir"; \
		(cd $$dir && go mod tidy); \
	done
	@echo "✅ Done tidying all Go modules."

fmt-all:
	@echo "🧹 Running go fmt in all modules..."
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "🎨 go fmt in $$dir"; \
		(cd $$dir && go fmt ./...); \
	done
	@echo "✅ Done formatting all Go modules."

vet-all:
	@echo "🕵️ Running go vet in all modules..."
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "🔍 go vet in $$dir"; \
		(cd $$dir && go vet ./...); \
	done
	@echo "✅ Done vetting all Go modules."

build-all:
	@echo "🏗️ Building all Go modules (without output)..."
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "🏗️ go build in $$dir"; \
		(cd $$dir && go build -o /dev/null ./...); \
	done
	@echo "✅ All modules compiled successfully."

update-all:
	@echo "⬆️ Updating only direct dependencies in all modules..."
	@find . -name "go.mod" -exec dirname {} \; | while read dir; do \
		echo "📦 Updating direct dependencies in $$dir"; \
		cd $$dir && \
		for dep in $$(go list -m -f '{{if not .Indirect}}{{.Path}}{{end}}' all); do \
			echo "   ↪️ go get -u $$dep"; \
			go get -u $$dep; \
		done; \
		go mod tidy; \
	done
	@echo "✅ All direct dependencies updated and tidied."


# Full code quality sweep
clean-all: tidy-all fmt-all vet-all

full: update-all tidy-all fmt-all vet-all build-all
