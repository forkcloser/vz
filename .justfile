# This file is the project's own — add recipes below. Keep the import: it
# mounts every shared limen task under `just do ...`.
import '.limen/just/main.just'

# The binding links Virtualization.framework: cgo is mandatory, and the module
# builds on darwin/arm64 only.
export GO_CGO := '1'

# The FIRST recipe defined here becomes `just`'s default.
lint: do::lint::go::default do::lint::go::deadcode do::lint::default
fix: do::fix::go::default do::fix::default
test: do::test::go::unit do::test::go::race
# The security workflow runs `just security`.
security: do::security::default
