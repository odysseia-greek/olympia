# Hippokrates

Integration/behavior tests for Odysseia services using `godog` (BDD `.feature` files) through Go `test`.

## What is currently tested

Current feature coverage is in `features/homeros.feature`.

Scenarios currently validate the Homeros GraphQL gateway:
- Gateway health endpoint is reachable.
- Grammar lookup returns expected declension rules for known words.
- Text flow works end-to-end: fetch text options, create a text, submit official translation, and expect perfect average Levenshtein.
- Word analysis returns complete analysis data.

All current scenarios are tagged with `@homeros`.

## Prerequisites

Run from `olympia/hippokrates`.

The tests call running services (they are not unit tests). In practice:
- `HOMEROS_SERVICE` should point to your gateway base URL (for example `http://localhost:8080`).
- Depending on your setup, `HERODOTOS_SERVICE` and `DIONYSIOS_SERVICE` may also be needed for backend calls used during scenarios.

Example:

```bash
export HOMEROS_SERVICE="http://localhost:8080"
export HERODOTOS_SERVICE="http://localhost:5001"
export DIONYSIOS_SERVICE="http://localhost:5002"
```

## Run all Hippokrates tests (terminal)

```bash
cd hippokrates
go test ./... -v
```

## Run tagged scenarios only (terminal)

`TestMain` reads a custom `-tags=` argument and forwards it to Godog tags.

Use `-args` so the tag reaches the test binary (instead of `go test` build tags):

```bash
cd hippokrates
go test ./... -v -args -tags=@homeros
```

You can also combine tags using Godog expressions, for example:

```bash
go test ./... -v -args '-tags=@homeros&&~@skip'
```

## Run tagged scenarios from GoLand (Program Arguments)

1. Open `Run | Edit Configurations...`.
2. Create or select a `Go Test` configuration for the `hippokrates` package/directory.
3. Set `Working directory` to `.../olympia/hippokrates`.
4. In `Program arguments`, set:

```text
-tags=@homeros -v
```

5. Add required environment variables in the run configuration (`HOMEROS_SERVICE`, and if needed `HERODOTOS_SERVICE` / `DIONYSIOS_SERVICE`).
6. Run the configuration.

## Notes on tags

- Tags come from `features/*.feature` files (for example `@homeros`).
- If no `-tags=` is provided, all scenarios in embedded feature files are executed.
