# Paths
SOKRATES_SCHEMA_SRC := ../../apologia/sokrates/graph/sokrates.graphqls
SOKRATES_SCHEMA_DEST := ./graph/sokrates.graphqls

SHARED_SCHEMA_SRC := ./graph/shared.graphqls   # Homeros' shared schema
SOKRATES_SHARED_SCHEMA_DEST := ../../apologia/sokrates/graph/shared.graphqls  # Copy to Sokrates
ALEXANDROS_SHARED_SCHEMA_DEST := ../../makedonia/alexandros/graph/shared.graphqls  # Copy to Alexandros

ALEXANDROS_SCHEMA_SRC := ../../makedonia/alexandros/graph/alexandros.graphqls
ALEXANDROS_SCHEMA_DEST := ./graph/alexandros.graphqls

HERODOTOS_SCHEMA_SRC := ../../ionia/herodotos/graph/schema.graphqls
HERODOTOS_SCHEMA_DEST := ./graph/herodotos.graphqls

# Colors for logging
BLUE := \033[1;34m
GREEN := \033[1;32m
YELLOW := \033[1;33m
NC := \033[0m  # No Color

GQLGEN = go tool gqlgen generate

# Ensure the graph directory exists
ensure_dirs:
	@mkdir -p $(dir $(SOKRATES_SCHEMA_DEST)) $(dir $(ALEXANDROS_SCHEMA_DEST)) $(dir $(HERODOTOS_SCHEMA_DEST))
	@echo "$(BLUE)[INFO]$(NC) Ensured destination directories exist."

# Copy Sokrates schema only if needed
sokrates_schema: ensure_dirs
	@if cmp -s $(SOKRATES_SCHEMA_SRC) $(SOKRATES_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in Sokrates schema, skipping copy."; \
	else \
		cp $(SOKRATES_SCHEMA_SRC) $(SOKRATES_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Sokrates schema copied successfully!"; \
	fi

# Copy Alexandros schema only if needed
alexandros_schema: ensure_dirs
	@if cmp -s $(ALEXANDROS_SCHEMA_SRC) $(ALEXANDROS_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in Alexandros schema, skipping copy."; \
	else \
		cp $(ALEXANDROS_SCHEMA_SRC) $(ALEXANDROS_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Alexandros schema copied successfully!"; \
	fi

# Import Herodotos as a schema extension. Two root fields are renamed locally
# because Homeros already exposes health and Alexandros already exposes text.
herodotos_schema: ensure_dirs
	@sed \
		-e 's/^type Query {/extend type Query {/' \
		-e 's/^  health: Health!/  herodotosHealth: Health!/' \
		-e 's/^  text(input:/  corpusText(input:/' \
		$(HERODOTOS_SCHEMA_SRC) > $(HERODOTOS_SCHEMA_DEST).tmp; \
	if cmp -s $(HERODOTOS_SCHEMA_DEST).tmp $(HERODOTOS_SCHEMA_DEST); then \
		rm $(HERODOTOS_SCHEMA_DEST).tmp; \
		echo "$(YELLOW)[SKIP]$(NC) No changes in Herodotos schema, skipping copy."; \
	else \
		mv $(HERODOTOS_SCHEMA_DEST).tmp $(HERODOTOS_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Herodotos schema copied successfully!"; \
	fi

# Copy Shared schema only if needed
shared_schema: ensure_dirs
	@if cmp -s $(SHARED_SCHEMA_SRC) $(SOKRATES_SHARED_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in shared schema, skipping copy."; \
	else \
		cp $(SHARED_SCHEMA_SRC) $(SOKRATES_SHARED_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Shared schema copied successfully!"; \
	fi

	@if cmp -s $(SHARED_SCHEMA_SRC) $(ALEXANDROS_SHARED_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in shared schema, skipping copy."; \
	else \
		cp $(SHARED_SCHEMA_SRC) $(ALEXANDROS_SHARED_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Shared schema copied successfully!"; \
	fi

.PHONY: gqlgen
gqlgen:
	@echo "🔄 Regenerating GraphQL code with gqlgen..."
	@$(GQLGEN)
	@echo "✅ GraphQL code regeneration complete!"

# Full sync: Runs both tasks
sync_schemas: shared_schema sokrates_schema alexandros_schema herodotos_schema
	@echo "$(BLUE)[INFO]$(NC) Schema synchronization completed."
