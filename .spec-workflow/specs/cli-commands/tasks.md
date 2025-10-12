# Tasks Document

## Implementation Tasks for cli-commands Spec

- [ ] 1. Initialize Go module and dependencies
  - Files: `go.mod`, `go.sum`
  - Initialize Go module with `go mod init github.com/veenone/mvnenv-win`
  - Add Cobra dependency: `github.com/spf13/cobra v1.8+`
  - Purpose: Establish project foundation and CLI framework dependency
  - _Leverage: None (foundation task)_
  - _Requirements: All (foundation)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in Go modules and dependency management | Task: Initialize Go module for mvnenv-win project and add Cobra CLI framework dependency v1.8+, ensuring proper module path and dependency versioning | Restrictions: Use Go 1.21+ compatible dependencies, do not add unnecessary dependencies, ensure reproducible builds | Success: go.mod created with correct module path, Cobra dependency added and working, `go mod tidy` runs without errors | Instructions: After completing, edit tasks.md and change this task from [ ] to [x]_

- [ ] 2. Create VERSION file for version management
  - Files: `VERSION`
  - Create VERSION file in project root containing version number (e.g., "1.0.0")
  - Document version file format (single line with semver version)
  - Purpose: Maintain single source of truth for application version
  - _Leverage: None (foundation task)_
  - _Requirements: Req 2 (Version Display Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: DevOps Engineer with expertise in version management and build automation | Task: Create VERSION file in project root containing initial version "1.0.0" following semantic versioning, establishing single source for version information that will be embedded at build time | Restrictions: VERSION file must contain only version number (no prefix, no newline at end), must follow semver format (MAJOR.MINOR.PATCH), file in project root | Success: VERSION file created with "1.0.0", file contains only version string, ready for build-time embedding | Instructions: After completing, edit tasks.md and change this task from [ ] to [x]_

- [ ] 3. Create main entry point with version embedding
  - Files: `cmd/mvnenv/main.go`
  - Create main.go with Cobra root command initialization
  - Add version variable for build-time embedding via ldflags
  - Implement fallback to read VERSION file during development
  - Setup root command with --version flag showing embedded/file version
  - Configure SilenceErrors and SilenceUsage for controlled error output
  - Purpose: Establish CLI application entry point with proper version management
  - _Leverage: Cobra command initialization patterns, Go build ldflags_
  - _Requirements: Req 1 (Command Framework), Req 2 (Version Display), Req 6 (Help System)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go CLI Developer specializing in Cobra framework and build-time variable injection | Task: Create main.go entry point implementing root command with Cobra and version embedding following Requirements 1, 2, and 6, using var Version string for ldflags embedding with fallback to read VERSION file, configuring proper error handling per design.md Component 1 | Restrictions: Must set SilenceErrors=true and SilenceUsage=true, version variable must be package-level var for ldflags, fallback reads VERSION file if Version is empty, ensure clean initialization pattern | Success: Application runs and displays help with `mvnenv --help`, `mvnenv --version` shows version from VERSION file during development, version embeddable via `go build -ldflags "-X main.Version=$(cat VERSION)"`, error handling configured correctly | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 4. Implement "commands" command
  - Files: `cmd/mvnenv/commands.go`
  - Create NewCommandsCmd() function returning cobra.Command
  - Implement runCommands() to list all available commands
  - Add command to root in main.go
  - Purpose: Provide command discovery per pyenv-win conventions
  - _Leverage: design.md Component 2 pattern_
  - _Requirements: Req 1 (Command Discovery)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in Cobra CLI framework | Task: Implement "commands" command following Requirement 1 and design.md Component 2, listing all available mvnenv commands (one per line, plain text) | Restrictions: Must follow design.md Component 2 implementation exactly, output plain text only (no colors/emojis), do not show hidden commands | Success: `mvnenv commands` lists all commands correctly, output is plain text list, follows pyenv-win format | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 5. Implement "version" command
  - Files: `cmd/mvnenv/version.go`
  - Create NewVersionCmd() function returning cobra.Command
  - Implement runVersion() with placeholder output
  - Add command to root in main.go
  - Purpose: Display currently active Maven version
  - _Leverage: design.md Component 3 pattern_
  - _Requirements: Req 2 (Version Display Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in CLI command implementation | Task: Implement "version" command following Requirement 2 and design.md Component 3, displaying current Maven version (placeholder: "3.9.4") or "No Maven version set" | Restrictions: Must follow design.md Component 3 exactly, use placeholder until core-version-management spec implemented, plain text output only | Success: `mvnenv version` outputs version number or "No Maven version set", plain text format, no errors | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 6. Implement "versions" command
  - Files: `cmd/mvnenv/versions.go`
  - Create NewVersionsCmd() function returning cobra.Command
  - Implement runVersions() to list installed versions with `*` for current
  - Add command to root in main.go
  - Purpose: List all installed Maven versions with current marked
  - _Leverage: design.md Component 4 pattern_
  - _Requirements: Req 2 (Version Display Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with CLI and formatting expertise | Task: Implement "versions" command following Requirement 2 and design.md Component 4, listing installed versions with current marked by `* ` prefix (placeholder versions) | Restrictions: Must follow design.md Component 4, use placeholder data, format as "* 3.9.4" for current and "  3.8.6" for others, plain text only | Success: `mvnenv versions` lists versions with proper formatting, current version marked with asterisk, plain text output | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 7. Implement "global" command
  - Files: `cmd/mvnenv/global.go`
  - Create NewGlobalCmd() function with MaximumNArgs(1)
  - Implement runGlobal() for both get and set operations
  - Set operations: no output on success (pyenv-win convention)
  - Add command to root in main.go
  - Purpose: Set or show global Maven version
  - _Leverage: design.md Component 5 pattern_
  - _Requirements: Req 3 (Version Selection Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in CLI argument parsing and command implementation | Task: Implement "global" command following Requirement 3 and design.md Component 5, handling both set and get operations with placeholder logic | Restrictions: Must follow design.md Component 5 pattern, no output on set success (pyenv-win convention), use cobra.MaximumNArgs(1), plain text only | Success: `mvnenv global` shows version, `mvnenv global 3.9.4` sets silently (no output), follows pyenv-win behavior | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 8. Implement "local" command
  - Files: `cmd/mvnenv/local.go`
  - Create NewLocalCmd() function with MaximumNArgs(1)
  - Implement runLocal() for both get and set operations
  - Set operations: no output on success
  - Add command to root in main.go
  - Purpose: Set or show local (project) Maven version
  - _Leverage: design.md Component 5 pattern (similar to global)_
  - _Requirements: Req 3 (Version Selection Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in file operations and CLI commands | Task: Implement "local" command following Requirement 3 and design.md Component 5, handling .maven-version file operations with placeholder logic | Restrictions: Must follow design.md Component 5 pattern, no output on set success, use cobra.MaximumNArgs(1), plain text only | Success: `mvnenv local` shows version, `mvnenv local 3.8.6` sets silently, follows pyenv-win conventions | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 9. Implement "shell" command
  - Files: `cmd/mvnenv/shell.go`
  - Create NewShellCmd() function with --unset flag
  - Implement runShell() for get/set/unset operations
  - Handle --unset flag for clearing shell version
  - Add command to root in main.go
  - Purpose: Set or show shell-specific Maven version
  - _Leverage: design.md Component 5 shell pattern_
  - _Requirements: Req 3 (Version Selection Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in environment variables and shell integration | Task: Implement "shell" command following Requirement 3 and design.md Component 5 shell pattern, handling get/set/unset operations with --unset flag | Restrictions: Must follow design.md shell command pattern, support --unset flag, no output on set success, plain text only | Success: `mvnenv shell` shows version, `mvnenv shell 3.6.3` sets silently, `mvnenv shell --unset` clears silently | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 10. Implement "install" command
  - Files: `cmd/mvnenv/install.go`
  - Create NewInstallCmd() with -l (list) and -q (quiet) flags
  - Implement runInstall() for listing and installation (placeholder logic)
  - Support multiple version arguments
  - Add command to root in main.go
  - Purpose: Install Maven versions or list available versions
  - _Leverage: design.md Component 6 pattern_
  - _Requirements: Req 4 (Installation and Listing Commands), Req 7 (Global Flags)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in CLI flags and command arguments | Task: Implement "install" command following Requirements 4 and 7 and design.md Component 6, supporting -l (list) and -q (quiet) flags with variadic version arguments | Restrictions: Must follow design.md Component 6 exactly, support multiple versions, respect -q flag for quiet output, use placeholder logic for actual installation | Success: `mvnenv install -l` lists versions, `mvnenv install 3.9.4` shows installation message, `mvnenv install -q 3.9.4` is silent, multiple versions supported | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 11. Implement "uninstall" command
  - Files: `cmd/mvnenv/uninstall.go`
  - Create NewUninstallCmd() requiring at least one argument
  - Implement runUninstall() supporting multiple versions
  - Use cobra.MinimumNArgs(1) for validation
  - Add command to root in main.go
  - Purpose: Uninstall one or more Maven versions
  - _Leverage: design.md Component 7 pattern_
  - _Requirements: Req 4 (Installation and Listing Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in command argument validation and error handling | Task: Implement "uninstall" command following Requirement 4 and design.md Component 7, requiring at least one version argument and supporting multiple versions | Restrictions: Must follow design.md Component 7 pattern, use cobra.MinimumNArgs(1), support multiple versions, use placeholder logic | Success: `mvnenv uninstall 3.9.4` works with message, `mvnenv uninstall 3.9.4 3.8.6` handles multiple versions, `mvnenv uninstall` shows error | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 12. Implement "update" command
  - Files: `cmd/mvnenv/update.go`
  - Create NewUpdateCmd() with no arguments
  - Implement runUpdate() to update version cache (placeholder)
  - Add command to root in main.go
  - Purpose: Update cached version database from repositories
  - _Leverage: design.md Component 8 update pattern_
  - _Requirements: Req 5 (Utility Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with CLI command implementation expertise | Task: Implement "update" command following Requirement 5 and design.md Component 8 update pattern, updating version cache with placeholder logic | Restrictions: Must follow design.md Component 8 pattern, no arguments required, simple output message, plain text only | Success: `mvnenv update` displays "Updating version cache..." message and completes successfully | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 13. Implement "rehash" command
  - Files: `cmd/mvnenv/rehash.go`
  - Create NewRehashCmd() with no arguments
  - Implement runRehash() to rebuild shims (placeholder)
  - Silent on success (no output)
  - Add command to root in main.go
  - Purpose: Rebuild shim executables
  - _Leverage: design.md Component 8 rehash pattern_
  - _Requirements: Req 5 (Utility Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with understanding of CLI conventions and silent operations | Task: Implement "rehash" command following Requirement 5 and design.md Component 8 rehash pattern, with no output on success (pyenv-win convention) | Restrictions: Must follow design.md Component 8 pattern, completely silent on success, no arguments, placeholder logic only | Success: `mvnenv rehash` completes silently with exit code 0, no output displayed | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 14. Implement "which" command
  - Files: `cmd/mvnenv/which.go`
  - Create NewWhichCmd() requiring exactly one argument
  - Implement runWhich() to show command path (placeholder)
  - Use cobra.ExactArgs(1) for validation
  - Add command to root in main.go
  - Purpose: Display full path to an executable in current Maven
  - _Leverage: design.md Component 8 which pattern_
  - _Requirements: Req 5 (Utility Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in argument validation and path operations | Task: Implement "which" command following Requirement 5 and design.md Component 8 which pattern, displaying executable path with placeholder logic | Restrictions: Must follow design.md Component 8 pattern, use cobra.ExactArgs(1), output full Windows path with backslashes, plain text only | Success: `mvnenv which mvn` displays path like "C:\\Users\\user\\.mvnenv\\versions\\3.9.4\\bin\\mvn.cmd", proper error if no argument | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 15. Implement "latest" command
  - Files: `cmd/mvnenv/latest.go`
  - Create NewLatestCmd() with optional prefix argument
  - Implement runLatest() to show latest version (placeholder)
  - Use cobra.MaximumNArgs(1) for optional prefix
  - Add command to root in main.go
  - Purpose: Show latest installed or known version matching optional prefix
  - _Leverage: design.md Component 8 latest pattern_
  - _Requirements: Req 5 (Utility Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in version comparison and filtering | Task: Implement "latest" command following Requirement 5 and design.md Component 8 latest pattern, showing latest version with optional prefix filter using placeholder logic | Restrictions: Must follow design.md Component 8 pattern, use cobra.MaximumNArgs(1), support optional prefix like "3.8", plain text output | Success: `mvnenv latest` shows "3.9.4", `mvnenv latest 3.8` shows "3.8.6" (placeholder), proper version filtering logic structure | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 16. Create error handling utilities
  - Files: `cmd/mvnenv/errors.go`
  - Implement formatError() function for consistent "Error: " prefix
  - Define common error variables (ErrVersionNotInstalled, etc.)
  - Create helper for stderr output formatting
  - Purpose: Provide consistent error handling across all commands
  - _Leverage: design.md Error Handling section_
  - _Requirements: Req 8 (Error Handling), Req 9 (Output Formatting)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in error handling patterns and consistent error messaging | Task: Create error handling utilities following Requirements 8 and 9 and design.md Error Handling section, providing consistent error formatting and common error types | Restrictions: Must follow design.md error patterns exactly, all errors must have "Error: " prefix, use stderr for errors, define exported error variables | Success: Error formatting functions work correctly, common errors defined, consistent "Error: " prefix, proper stderr output | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 17. Create output utilities
  - Files: `cmd/mvnenv/output.go`
  - Implement shouldPrint() helper for --quiet flag checking
  - Create printVersion(), printVersionList() formatters
  - Add Windows path formatting helper (backslashes)
  - Purpose: Provide consistent output formatting across commands
  - _Leverage: design.md Output Formatting section_
  - _Requirements: Req 7 (Global Flags), Req 9 (Output Formatting)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in output formatting and cross-platform path handling | Task: Create output utilities following Requirements 7 and 9 and design.md Output Formatting section, providing consistent plain-text output helpers | Restrictions: Must follow design.md output patterns, plain text only (no colors/emojis), respect --quiet flag, Windows backslash paths | Success: Output helpers work correctly, --quiet flag respected, version formatting correct, Windows paths use backslashes | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 18. Add unit tests for root command
  - Files: `cmd/mvnenv/main_test.go`
  - Test root command initialization
  - Test --help and --version flags
  - Test error handling and exit codes
  - Verify all commands registered
  - Purpose: Ensure root command and framework work correctly
  - _Leverage: Go testing package, design.md Testing Strategy_
  - _Requirements: All (foundation testing)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in Go testing and CLI testing patterns | Task: Create comprehensive unit tests for root command following design.md Testing Strategy, testing initialization, help, version, and command registration | Restrictions: Must test command structure not business logic, use table-driven tests, ensure >90% coverage of main.go, test exit codes | Success: All root command tests pass, >90% coverage, help/version flags tested, all commands verified registered | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 19. Add unit tests for info commands
  - Files: `cmd/mvnenv/commands_test.go`, `cmd/mvnenv/version_test.go`, `cmd/mvnenv/versions_test.go`
  - Test commands, version, versions command execution
  - Test output formatting and content
  - Test error scenarios
  - Verify plain text output (no colors/emojis)
  - Purpose: Ensure info commands work correctly
  - _Leverage: Go testing package, design.md Testing Strategy_
  - _Requirements: Req 1 (Command Discovery), Req 2 (Version Display)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer specializing in CLI output testing and Go test frameworks | Task: Create unit tests for commands, version, versions commands following design.md Testing Strategy and Requirements 1-2, testing output format and content | Restrictions: Must use table-driven tests, verify plain text output only, test placeholder outputs, ensure >90% coverage per file | Success: All info command tests pass, output format verified, plain text confirmed, >90% coverage, no colors/emojis in output | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 20. Add unit tests for version selection commands
  - Files: `cmd/mvnenv/global_test.go`, `cmd/mvnenv/local_test.go`, `cmd/mvnenv/shell_test.go`
  - Test global, local, shell command execution
  - Test both get and set operations
  - Test shell --unset flag
  - Verify no output on successful set operations
  - Purpose: Ensure version selection commands work correctly
  - _Leverage: Go testing package, design.md Testing Strategy_
  - _Requirements: Req 3 (Version Selection Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in command argument testing and silent operation verification | Task: Create unit tests for global, local, shell commands following design.md Testing Strategy and Requirement 3, testing get/set operations and silent success | Restrictions: Must verify no output on set success (pyenv-win convention), test both with and without arguments, test shell --unset flag, >90% coverage | Success: All version selection tests pass, silent set operations verified, get operations show output, shell --unset tested, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 21. Add unit tests for install/uninstall commands
  - Files: `cmd/mvnenv/install_test.go`, `cmd/mvnenv/uninstall_test.go`
  - Test install command with -l and -q flags
  - Test multiple version arguments
  - Test uninstall with multiple versions
  - Test argument validation (minimum args)
  - Purpose: Ensure install/uninstall commands work correctly
  - _Leverage: Go testing package, design.md Testing Strategy_
  - _Requirements: Req 4 (Installation Commands), Req 7 (Global Flags)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer specializing in flag testing and variadic argument validation | Task: Create unit tests for install and uninstall commands following design.md Testing Strategy and Requirements 4 and 7, testing flags and multiple arguments | Restrictions: Must test -l and -q flags independently and together, verify multiple version support, test argument validation, >90% coverage | Success: All install/uninstall tests pass, flags tested thoroughly, multiple versions verified, quiet mode confirmed, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 22. Add unit tests for utility commands
  - Files: `cmd/mvnenv/update_test.go`, `cmd/mvnenv/rehash_test.go`, `cmd/mvnenv/which_test.go`, `cmd/mvnenv/latest_test.go`
  - Test update, rehash, which, latest commands
  - Test argument validation (exact, maximum args)
  - Test rehash silent operation
  - Test which and latest with optional/required args
  - Purpose: Ensure utility commands work correctly
  - _Leverage: Go testing package, design.md Testing Strategy_
  - _Requirements: Req 5 (Utility Commands)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in argument validation testing and command behavior verification | Task: Create unit tests for update, rehash, which, latest commands following design.md Testing Strategy and Requirement 5, testing all argument patterns | Restrictions: Must verify rehash silence, test which cobra.ExactArgs(1), test latest cobra.MaximumNArgs(1), test update no args, >90% coverage | Success: All utility command tests pass, argument validation verified, rehash silence confirmed, which/latest args tested, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 23. Add unit tests for error and output utilities
  - Files: `cmd/mvnenv/errors_test.go`, `cmd/mvnenv/output_test.go`
  - Test formatError() function
  - Test shouldPrint() with --quiet flag
  - Test output formatting helpers
  - Test Windows path formatting
  - Purpose: Ensure utility functions work correctly
  - _Leverage: Go testing package, design.md Testing Strategy_
  - _Requirements: Req 8 (Error Handling), Req 9 (Output Formatting)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer specializing in utility function testing and edge case coverage | Task: Create unit tests for error and output utilities following design.md Testing Strategy and Requirements 8-9, testing formatting and flag handling | Restrictions: Must test error prefix consistency, verify quiet flag behavior, test Windows backslash paths, ensure 100% coverage of utilities | Success: All utility tests pass, error formatting verified, quiet flag tested, path formatting confirmed, 100% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 24. Add integration tests for command execution
  - Files: `cmd/mvnenv/integration_test.go`
  - Test full command execution flow
  - Test flag combinations
  - Test error output goes to stderr
  - Test exit codes
  - Purpose: Ensure commands work together correctly
  - _Leverage: Go testing package, design.md Integration Testing section_
  - _Requirements: All command requirements_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Integration Engineer with expertise in end-to-end CLI testing and Go integration tests | Task: Create integration tests for command execution flow following design.md Integration Testing section, testing complete command scenarios | Restrictions: Must test real command execution paths, verify stderr for errors, test exit codes, test flag combinations, use testable examples | Success: Integration tests pass covering main scenarios, stderr verified for errors, exit codes confirmed, flag combinations tested | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 25. Add performance benchmarks
  - Files: `cmd/mvnenv/benchmark_test.go`
  - Create benchmarks for version, versions, help commands
  - Verify <50ms for version/--version
  - Verify <100ms for help and versions
  - Track command startup overhead
  - Purpose: Ensure performance requirements are met
  - _Leverage: Go benchmark testing, design.md Performance Testing section_
  - _Requirements: Req 10 (Command Execution Performance)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Performance Engineer with expertise in Go benchmarking and CLI performance optimization | Task: Create performance benchmarks following design.md Performance Testing section and Requirement 10, verifying performance targets | Restrictions: Must benchmark critical commands only, verify <50ms for version, <100ms for help/versions, track startup overhead separately | Success: Benchmarks run successfully, performance targets met or documented gaps, startup overhead measured | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 26. Create end-to-end test scenarios
  - Files: `test/e2e/cli_test.go`
  - Test first-time user scenario
  - Test command discovery flow
  - Test help system navigation
  - Test error recovery scenarios
  - Purpose: Validate complete user workflows
  - _Leverage: Go testing package, design.md End-to-End Testing section_
  - _Requirements: All requirements (user experience validation)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Automation Engineer specializing in user acceptance testing and E2E test design | Task: Create end-to-end test scenarios following design.md End-to-End Testing section, testing complete user workflows and help discovery | Restrictions: Must test realistic user scenarios, verify help system usability, test error recovery, ensure tests are maintainable | Success: E2E tests pass covering critical user journeys, help discovery works, error recovery verified, tests are reliable | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 27. Add command help text and examples
  - Files: All command files (`cmd/mvnenv/*.go`)
  - Add comprehensive help text to all commands
  - Include usage examples in Long descriptions
  - Add flag descriptions
  - Ensure consistency across commands
  - Purpose: Provide helpful documentation in CLI
  - _Leverage: Cobra help system, design.md Help System section_
  - _Requirements: Req 6 (Help System)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer with expertise in CLI documentation and user assistance | Task: Add comprehensive help text and examples to all commands following Requirement 6 and design.md Help System section, ensuring clarity and consistency | Restrictions: Must include usage, short, long descriptions with examples, describe all flags, follow pyenv-win help style, plain text only | Success: All commands have complete help text, examples included, --help output is clear and helpful, consistent style across commands | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 28. Create README documentation
  - Files: `README.md`
  - Document all commands with examples
  - Add installation instructions (placeholder)
  - Describe pyenv-win compatibility
  - Include troubleshooting section
  - Purpose: Provide user-facing documentation
  - _Leverage: requirements.md command reference, design.md_
  - _Requirements: All requirements (documentation)_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Technical Writer specializing in open-source project documentation and CLI tools | Task: Create comprehensive README documenting all CLI commands, installation, and usage following requirements.md command reference and pyenv-win compatibility | Restrictions: Must document all commands from requirements.md Command Reference section, note pyenv-win compatibility, include examples for each command, plain markdown | Success: README is complete and clear, all commands documented with examples, installation section present, pyenv-win compatibility noted | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 29. Final integration and code review
  - Files: All files in `cmd/mvnenv/`
  - Review all command implementations for consistency
  - Verify error handling follows patterns
  - Ensure all commands follow pyenv-win conventions
  - Run all tests and verify >90% coverage
  - Check performance benchmarks
  - Purpose: Ensure specification is fully implemented
  - _Leverage: All previous tasks, requirements.md Success Criteria_
  - _Requirements: All requirements_
  - _Prompt: Implement the task for spec cli-commands, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Senior Go Developer and Code Reviewer with expertise in CLI applications and code quality | Task: Perform final integration review of all CLI commands, verifying consistency, error handling, pyenv-win conventions, test coverage, and performance per requirements.md Success Criteria | Restrictions: Must verify all Success Criteria from requirements.md are met, ensure >90% test coverage, confirm pyenv-win alignment, validate plain text output | Success: All Success Criteria from requirements.md verified and met, tests pass with >90% coverage, code is consistent and clean, ready for next spec integration | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

## Task Execution Notes

### VERSION File and Build Process

**Development Mode:**
- VERSION file in project root contains version string (e.g., "1.0.0")
- Application reads VERSION file if embedded version is empty
- Allows version changes without recompilation during development

**Production Build:**
- Version embedded at build time using ldflags
- Build command: `go build -ldflags "-X main.Version=$(cat VERSION)" -o mvnenv.exe cmd/mvnenv/main.go`
- Embedded version takes precedence; VERSION file not needed in distribution
- Single binary with no external dependencies

**Version Management:**
- Update VERSION file to change version
- CI/CD should embed version automatically during build
- Binary shows embedded version via `mvnenv --version`

### Parallel Execution Opportunities

Tasks can be executed in parallel within these groups:

**Group 1: Foundation (Tasks 1-3)**
- Task 1 (go modules) must complete first
- Tasks 2 (VERSION file) and 3 (main.go) can be done in parallel after task 1, but task 3 depends on task 2 for reading VERSION during development

**Group 2: Command Implementations (Tasks 4-15)**
- Commands 4-15 can be implemented independently after task 3 is complete
- Each command is self-contained in its own file
- All follow the same pattern from design.md

**Group 3: Utility Files (Tasks 16-17)**
- Can be implemented in parallel after understanding error/output patterns
- Both are independent helper modules

**Group 4: Unit Tests (Tasks 18-23)**
- Can be written in parallel after corresponding commands are implemented
- Each test file is independent

**Group 5: Integration and Documentation (Tasks 24-28)**
- Can proceed in parallel after all implementations complete
- Each focuses on different aspect (integration, performance, E2E, docs)

### Sequential Dependencies

- Task 1 must complete before task 2 and 3
- Task 2 (VERSION file) should complete before task 3 (main.go needs it for fallback)
- Task 3 must complete before tasks 4-15
- Each command (tasks 4-15) must complete before its unit test
- All implementations (tasks 4-17) must complete before integration tests (24-26)
- Task 29 must be last (final review)

### Placeholder Pattern

Business logic for version management, installation, etc. will be implemented in other specs:
- **core-version-management**: Install, uninstall, version resolution
- **shim-system-implementation**: Rehash, which command logic
- **nexus-repository-integration**: Repository management for install -l

For this spec, use placeholder outputs and TODO comments marking integration points.
