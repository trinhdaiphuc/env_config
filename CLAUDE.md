# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...
go test -v -race -covermode=atomic ./...   # what CI runs (.github/workflows/go.yml, Go 1.23.1)
go test -run TestLoadConfig ./             # single test
mkdir -p coverage && make test             # coverage + opens HTML report; make test fails without coverage/
```

`examples/` is a **separate module** — `cd examples && go run example.go` (it depends on the
published `env_config` version, not the working tree; add a `replace` directive to test local changes).

## Architecture

Single package `env_config`. One public entry point: `LoadConfig(cfg any) error` in `env.go`.
Note the README says `Load` — the actual exported name is `LoadConfig`.

The load path is a four-stage chain; read all four files before editing any one of them:

1. **`structs.go`** — `NewStruct` reflects over the struct, skipping fields without an `env` tag.
   Splits the tag on `;` into key + options, prefixes the key via `combineKeyPrefix`
   (nested struct keys become `PARENT_CHILD`), then asks the handler factory for a handler.
   Produces a tree of `Item` (`FieldItem` leaf / `StructItem` branch); `Load()` walks it.
   `FieldItem.Load()` reads `os.Getenv(key)` and dispatches via `lookupStrategy`:
   `complexTypeStrategies[type]` first, then `buildInTypeStrategies[kind]`; **an unsupported type
   is an error, not a no-op**. Both registries are guarded by `_strategyMu` because
   `RegisterStrategy` is exported.

   **Pointer fields are allocated lazily.** A nil pointer stays nil unless some leaf under it
   resolves to a value (`hasEnvValue`: env var present, or a non-empty `default=`). A nil pointer
   *section* is built against a throwaway `reflect.New` (`StructItem.alloc`) and only assigned to
   the real field (`StructItem.target`) once that check passes — so `cfg.Section != nil` means the
   section was actually configured.
2. **`type_handler.go`** — `TypeHandlerFactory` maps `reflect.Type` → handler. `Handle` returns
   `(Item, error)`; never return a nil `Item`, `StructItem.Load` would deref it. `time.Time` is
   explicitly registered so it is treated as a leaf, not recursed into as a struct. Any other
   struct (or pointer-to-struct) gets `StructHandler`, which recurses via `NewStruct`.
3. **`types.go`** — `TypeStrategy` per concrete type / kind, registered in the two maps in `init()`.
   **New supported types are added here**, never by branching in `structs.go`. Third parties use
   `RegisterStrategy(reflect.Type, TypeStrategy)`. Slice/int/uint/float strategies are generic over
   the constraints in `generic.go`.
4. **`tags.go`** — tag options (`default=`, `delimiter=`) are builders producing a *chain* sorted by
   `Priority()` (default=0 runs before delimiter=1), so the default value is substituted before it is
   split. `Apply` returns `string` or `[]string`; `parseOptionValue`/`parseOptionValues` in `types.go`
   narrow it back.

### Conventions worth keeping

- Empty env value + no default → strategies return `nil` and leave the zero value; they do not error.
  Malformed values do: parse errors (including per-element slice errors and integer overflow, via
  `field.Type().Bits()`) are wrapped and returned, never swallowed.
- Every strategy type-asserts its field kind and returns `fmt.Errorf("invalid type, expected ...")` on mismatch.
- `var _ Item = FieldItem{}` style compile-time interface checks at the top of the file.
- Tests are table-driven with `t.Run` subtests and `testify`; env vars set with `t.Setenv`/`os.Setenv`.
  Adding a type means adding a case to `types_test.go`.
