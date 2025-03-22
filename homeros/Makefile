# Paths
SOKRATES_SCHEMA_SRC := ../../apologia/sokrates/graph/sokrates.graphqls
SOKRATES_SCHEMA_DEST := ./graph/sokrates.graphqls

SHARED_SCHEMA_SRC := ./graph/shared.graphqls   # Homeros' shared schema
SHARED_SCHEMA_DEST := ../../apologia/sokrates/graph/shared.graphqls  # Copy to Sokrates

# Colors for logging
BLUE := \033[1;34m
GREEN := \033[1;32m
YELLOW := \033[1;33m
NC := \033[0m  # No Color

# Ensure the graph directory exists
ensure_dirs:
	@mkdir -p $(dir $(SOKRATES_SCHEMA_DEST)) $(dir $(SHARED_SCHEMA_DEST))
	@echo "$(BLUE)[INFO]$(NC) Ensured destination directories exist."

# Copy Sokrates schema only if needed
sokrates_schema: ensure_dirs
	@if cmp -s $(SOKRATES_SCHEMA_SRC) $(SOKRATES_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in Sokrates schema, skipping copy."; \
	else \
		cp $(SOKRATES_SCHEMA_SRC) $(SOKRATES_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Sokrates schema copied successfully!"; \
	fi

# Copy Shared schema only if needed
shared_schema: ensure_dirs
	@if cmp -s $(SHARED_SCHEMA_SRC) $(SHARED_SCHEMA_DEST); then \
		echo "$(YELLOW)[SKIP]$(NC) No changes in shared schema, skipping copy."; \
	else \
		cp $(SHARED_SCHEMA_SRC) $(SHARED_SCHEMA_DEST); \
		echo "$(GREEN)[UPDATED]$(NC) Shared schema copied successfully!"; \
	fi

# Full sync: Runs both tasks
sync_schemas: shared_schema sokrates_schema
	@echo "$(BLUE)[INFO]$(NC) Schema synchronization completed."
