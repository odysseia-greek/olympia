# Olympia

Olympia is the main repository for Odysseia-Greek related APIs and supporting services.

In less modest terms: this is the polis where the services gather, argue, and (usually) agree on Ancient Greek learning workflows.

## What Olympia contains

Each top-level folder is a focused component in the platform.

- `herodotos`: backend API for text retrieval and translation workflows.
- `homeros`: GraphQL gateway that exposes and combines backend functionality for clients.
- `pheidias`: frontend application (Vue/Vite) for learners.
- `hippokrates`: integration/behavior test suite (Godog + Go test) validating service health and end-to-end flows.
- `herakleitos`: indexing/seeding job that loads text corpora (embedded `rhema`) into search/index infrastructure.
- `melissos`: background processing job for dictionary/data handling and completion signaling.
- `protagoras`: seeding job for bootstrapping required data sets.

## Repository role in Odysseia-Greek

Olympia is the application-and-jobs repository.
- It contains runtime APIs (`herodotos`, `homeros`), UI (`pheidias`), tests (`hippokrates`), and operational data jobs (`herakleitos`, `melissos`, `protagoras`).
- Supporting libraries and sibling services live in other Odysseia-Greek repositories.

## Deployment

Deployment configuration and rollout flow are handled in **Mykenai**:
- https://github.com/odysseia-greek/mykenai

## Monorepo maintenance helpers

From repo root:

```bash
make tidy-all
make fmt-all
make vet-all
make build-all
```

Or run a full sweep:

```bash
make full
```
