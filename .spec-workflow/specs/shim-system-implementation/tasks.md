# Tasks Document

This document lists all implementation tasks for the Shim System Implementation feature. Each task includes detailed prompts for AI-assisted development.

## Task List

### Foundation Tasks

#### Task 1: Create Shim Executable Entry Point
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `cmd/shim/main.go`

**Description:** Create the main entry point for the shim executable that detects the command name from the executable filename and routes to the ShimExecutor.

**Dependencies:** None

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are implementing the shim executable entry point for mvnenv-win.

**Role:** Senior Go developer with expertise in command-line tools and process execution.

**Task:** Create cmd/shim/main.go with:

1. main() function:
   - Detect command name from os.Executable() path
   - Extract base filename and remove .exe extension
   - Create VersionResolver from core-version-management
   - Create ShimExecutor with resolver
   - Call executor.Execute(command, os.Args[1:])
   - Exit with returned exit code

2. detectCommand() string:
   - Get executable path with os.Executable()
   - Extract basename with filepath.Base()
   - Remove .exe extension if present with strings.TrimSuffix()
   - Return command name (e.g., "mvn", "mvnDebug")
   - Default to "mvn" on error

3. Error handling:
   - Print errors to stderr
   - Exit with code 1 on resolution errors
   - Exit with Maven's code on execution errors

4. Minimal imports:
   - Standard library only (fmt, os, path/filepath, strings)
   - Internal packages: internal/shim, internal/version

**Restrictions:**
- No external dependencies
- Keep main.go minimal (<50 lines)
- All logic in internal/shim package
- Never buffer or modify Maven output

**Success Criteria:**
- Executable detects command from filename
- Errors print to stderr
- Exit codes preserved from Maven
- Clean, readable code

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 2: Create ShimExecutor for Command Execution
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/executor.go`

**Description:** Implement ShimExecutor that handles version resolution, Maven path construction, process execution, and I/O forwarding.

**Dependencies:** Task 1, core-version-management spec completed

**Estimated Effort:** 4 hours

**_Prompt:**
```
You are implementing ShimExecutor for mvnenv-win's shim system.

**Role:** Senior Go developer with expertise in process management and I/O handling.

**Task:** Create internal/shim/executor.go with:

1. ShimExecutor struct:
   - resolver *version.VersionResolver
   - debug bool (from MVNENV_DEBUG env var)

2. NewShimExecutor(resolver *version.VersionResolver) *ShimExecutor:
   - Check os.Getenv("MVNENV_DEBUG") == "1"
   - Return executor with resolver and debug flag

3. Execute(command string, args []string) (int, error):
   - Record start time for performance tracking
   - Call resolver.ResolveVersion()
   - If error: return formatResolutionError()
   - Construct Maven path: filepath.Join(resolved.Path, "bin", command+".cmd")
   - Verify Maven binary exists with os.Stat()
   - If debug: log resolution details to stderr
   - Call executeMaven()
   - If debug: log total execution time
   - Return exit code

4. executeMaven(mavenPath string, args []string, resolved *version.ResolvedVersion) (int, error):
   - Create context.WithCancel() for signal handling
   - Create exec.CommandContext(ctx, mavenPath, args...)
   - Set MAVEN_HOME env var: append to os.Environ()
   - Forward I/O: cmd.Stdin/Stdout/Stderr = os.Stdin/Stdout/Stderr
   - Set working directory: cmd.Dir, _ = os.Getwd()
   - Setup signal handler for Ctrl+C: signal.Notify() -> cancel()
   - Start process: cmd.Start()
   - Wait for completion: cmd.Wait()
   - Extract exit code from *exec.ExitError
   - Return exit code (0 on success)

5. formatResolutionError(err error) error:
   - Check error type with version.IsVersionNotInstalledError()
   - Format user-friendly message with installation suggestion
   - Check version.IsNoVersionSetError()
   - Format message with global version suggestion
   - Default: wrap error with context

6. logDebug(...):
   - Output to stderr with [mvnenv] prefix
   - Include: command, args, version, source, path, MAVEN_HOME, timing

**Restrictions:**
- Use os/exec for process execution
- Use context for cancellation
- No buffering of I/O streams
- Direct forwarding: cmd.Stdin = os.Stdin
- Handle signals properly (Ctrl+C)
- Extract exact exit code from Maven

**Success Criteria:**
- Version resolution integrated correctly
- Process execution with I/O forwarding works
- Exit codes preserved exactly
- Signal handling (Ctrl+C) works
- Debug mode provides useful information
- Performance: <50ms overhead

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 3: Create ShimGenerator for Shim File Generation
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/generator.go`

**Description:** Implement ShimGenerator that creates shim executables by copying shim.exe and generates .cmd batch files.

**Dependencies:** None

**Estimated Effort:** 4 hours

**_Prompt:**
```
You are implementing ShimGenerator for mvnenv-win.

**Role:** Senior Go developer with file system operations expertise.

**Task:** Create internal/shim/generator.go with:

1. ShimGenerator struct:
   - shimsDir string (%USERPROFILE%\.mvnenv\shims)
   - shimBinary string (path to shim.exe)
   - versionsDir string (%USERPROFILE%\.mvnenv\versions)

2. NewShimGenerator(mvnenvRoot string) *ShimGenerator:
   - Initialize paths using filepath.Join()
   - shimsDir: filepath.Join(mvnenvRoot, "shims")
   - shimBinary: filepath.Join(mvnenvRoot, "bin", "shim.exe")
   - versionsDir: filepath.Join(mvnenvRoot, "versions")

3. GenerateShims() ([]string, error):
   - Create shims directory: os.MkdirAll(shimsDir, 0755)
   - Define core commands: []string{"mvn", "mvnDebug"}
   - Call discoverAdditionalCommands() for mvnyjp, etc.
   - Append additional commands to core commands
   - For each command:
     * Call generateShimFile(cmd, ".exe")
     * Call generateBatchShim(cmd)
     * Collect generated paths
   - Return list of generated paths

4. generateShimFile(command string, ext string) (string, error):
   - Construct destination: filepath.Join(shimsDir, command+ext)
   - Copy shim.exe to destination with copyFile()
   - Verify executable with verifyExecutable()
   - Return destination path

5. generateBatchShim(command string) (string, error):
   - Construct destination: filepath.Join(shimsDir, command+".cmd")
   - Create batch script:
     @echo off
     "%~dp0{command}.exe" %*
     exit /b %ERRORLEVEL%
   - Write with os.WriteFile(dest, []byte(script), 0755)
   - Return destination path

6. discoverAdditionalCommands() ([]string, error):
   - Read versionsDir with os.ReadDir()
   - For each version directory:
     * Read bin/ subdirectory
     * Look for .cmd files
     * Exclude mvn.cmd and mvnDebug.cmd
     * Collect unique command names in set
   - Return list of additional commands

7. copyFile(src, dst string) error:
   - Open source file
   - Create temp destination: dst+".tmp"
   - io.Copy() from source to temp
   - Close both files
   - Atomic rename: os.Rename(tmp, dst)
   - Clean up temp file on error

8. verifyExecutable(path string) error:
   - os.Stat() to check file exists
   - Verify not a directory
   - Check .exe or .cmd extension

**Restrictions:**
- Use atomic file operations (temp + rename)
- Handle missing versionsDir gracefully (no versions yet)
- Log warnings for non-critical errors
- Create directories automatically
- Windows-specific: .exe and .cmd files

**Success Criteria:**
- Shims generated for mvn and mvnDebug
- Additional commands discovered from installed versions
- Atomic generation (no partial files)
- Both .exe and .cmd variants created
- Existing shims overwritten cleanly

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 4: Create PathManager for Windows PATH Management
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/path_windows.go`

**Description:** Implement Windows registry-based PATH management for adding/removing shims directory from user PATH.

**Dependencies:** None

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are implementing Windows PATH management for mvnenv-win.

**Role:** Senior Go developer with Windows platform expertise.

**Task:** Create internal/shim/path_windows.go with:

1. Build tag: //go:build windows

2. PathManager struct:
   - shimsDir string

3. NewPathManager(mvnenvRoot string) *PathManager:
   - Initialize shimsDir: filepath.Join(mvnenvRoot, "shims")

4. AddShimsDirToPath() (bool, error):
   - Open registry key: registry.OpenKey(registry.CURRENT_USER, `Environment`, QUERY_VALUE|SET_VALUE)
   - Read current PATH: key.GetStringValue("Path")
   - Check if shimsDir already in PATH with isInPath() (case-insensitive)
   - If already present: return false, nil
   - Prepend shimsDir: newPath = shimsDir + ";" + currentPath
   - Write new PATH: key.SetStringValue("Path", newPath)
   - Broadcast environment change: broadcastEnvironmentChange()
   - Return true, nil

5. RemoveShimsDirFromPath() (bool, error):
   - Open registry key
   - Read current PATH
   - Remove shimsDir with removeFromPath()
   - If not present: return false, nil
   - Write new PATH
   - Broadcast environment change
   - Return true, nil

6. isInPath(path string, dir string) bool:
   - Split PATH by semicolon: strings.Split(path, ";")
   - Convert dir to lowercase for comparison
   - For each path element:
     * Trim whitespace
     * Compare lowercase
     * Return true if match
   - Return false

7. removeFromPath(path string, dir string) string:
   - Split PATH by semicolon
   - Filter out directory (case-insensitive)
   - Join remaining paths with semicolon

8. broadcastEnvironmentChange():
   - Use syscall to SendMessageTimeoutW
   - Broadcast WM_SETTINGCHANGE message
   - Makes PATH change visible to new processes

**Restrictions:**
- Windows-only code (build tag)
- Use golang.org/x/sys/windows/registry
- Case-insensitive PATH comparison
- Prepend to PATH (highest priority)
- User-level registry only (CURRENT_USER)
- Handle missing PATH value gracefully

**Success Criteria:**
- Shims directory added to user PATH
- PATH change visible in new terminal sessions
- Duplicate additions prevented
- Case-insensitive comparison works
- Removal works correctly

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Integration Tasks

#### Task 5: Integrate ShimGenerator with Version Installer
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Modify:**
- `internal/version/installer.go`

**Description:** Add automatic shim regeneration after Maven version installation.

**Dependencies:** Tasks 3, core-version-management spec completed

**Estimated Effort:** 1 hour

**_Prompt:**
```
You are integrating automatic shim regeneration into mvnenv-win's version installer.

**Role:** Senior Go developer with integration expertise.

**Task:** Update internal/version/installer.go:

1. Add field to VersionInstaller struct:
   - shimGenerator *shim.ShimGenerator

2. Update NewVersionInstaller():
   - Create ShimGenerator: shim.NewShimGenerator(mvnenvRoot)
   - Assign to installer.shimGenerator

3. Update InstallVersion() method:
   - After successful installation (extraction and verification)
   - Call installer.regenerateShims()
   - If regeneration fails: log warning, don't fail installation

4. Add regenerateShims() method:
   - Call shimGenerator.GenerateShims()
   - Log success: "Shims regenerated (%d files)"
   - Return error if generation fails

5. Similarly update UninstallVersion() method:
   - After successful uninstallation
   - Call installer.regenerateShims()
   - Log warning on failure

**Restrictions:**
- Don't fail installation if shim regeneration fails
- Log warnings for shim errors
- Provide suggestion to run 'mvnenv rehash' manually

**Success Criteria:**
- Shims automatically regenerated after install
- Shims automatically regenerated after uninstall
- Installation doesn't fail if regeneration fails
- Clear warning messages

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 6: Implement rehash CLI Command
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `cmd/mvnenv/cmd/rehash.go`

**Description:** Implement the `mvnenv rehash` command that manually triggers shim regeneration.

**Dependencies:** Task 3, cli-commands spec completed

**Estimated Effort:** 1 hour

**_Prompt:**
```
You are implementing the rehash command for mvnenv-win.

**Role:** Senior Go developer with Cobra CLI expertise.

**Task:** Create cmd/mvnenv/cmd/rehash.go:

1. Create Cobra command:
   - Use: "rehash"
   - Short: "Regenerate shim executables"
   - Long: Detailed description of what rehash does and when to use it

2. runRehash(cmd *cobra.Command, args []string) error:
   - Get mvnenvRoot from config or environment
   - Create ShimGenerator: shim.NewShimGenerator(mvnenvRoot)
   - Call generator.GenerateShims()
   - If error: return formatted error
   - Print success message: "Shims regenerated successfully (%d files)"
   - List generated shim names (without paths)

3. init() function:
   - Register command: rootCmd.AddCommand(rehashCmd)

4. Output format (plain text, no emojis):
   Regenerating shims...
   Shims regenerated successfully (6 files)
   - mvn
   - mvnDebug
   - mvnyjp

**Restrictions:**
- Follow pyenv-win output style (plain text)
- No emojis or decorative elements
- Clear error messages
- Show list of regenerated commands

**Success Criteria:**
- Command runs without errors
- Shims regenerated correctly
- Output shows what was generated
- Error messages clear and actionable

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Testing Tasks

#### Task 7: Unit Tests for ShimExecutor
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/executor_test.go`

**Description:** Comprehensive unit tests for ShimExecutor including resolution, execution, and error handling.

**Dependencies:** Task 2

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are writing unit tests for ShimExecutor in mvnenv-win.

**Role:** Senior Go developer with testing expertise.

**Task:** Create internal/shim/executor_test.go:

1. TestShimExecutor_Execute_Success:
   - Mock VersionResolver to return valid version
   - Create fake Maven binary (batch script that echoes args)
   - Call executor.Execute()
   - Verify: exit code 0
   - Verify: no error

2. TestShimExecutor_Execute_VersionNotInstalled:
   - Mock VersionResolver to return ErrVersionNotInstalled
   - Call executor.Execute()
   - Verify: exit code 1
   - Verify: error message suggests installation command

3. TestShimExecutor_Execute_NoVersionSet:
   - Mock VersionResolver to return ErrNoVersionSet
   - Call executor.Execute()
   - Verify: exit code 1
   - Verify: error message suggests setting global version

4. TestShimExecutor_Execute_MavenBinaryNotFound:
   - Mock VersionResolver to return valid but non-existent path
   - Call executor.Execute()
   - Verify: exit code 1
   - Verify: error mentions binary not found

5. TestShimExecutor_Execute_ExitCodePreserved:
   - Mock Maven binary that exits with code 42
   - Call executor.Execute()
   - Verify: exit code 42 returned

6. TestShimExecutor_Execute_ArgumentPassthrough:
   - Mock Maven binary that echoes arguments
   - Call executor.Execute("mvn", []string{"clean", "install"})
   - Capture output
   - Verify: arguments passed unchanged

7. TestShimExecutor_Execute_DebugMode:
   - Set MVNENV_DEBUG=1
   - Capture stderr
   - Call executor.Execute()
   - Verify: debug information logged to stderr
   - Verify: includes version, source, path, timing

**Restrictions:**
- Use testify/assert for assertions
- Mock VersionResolver with testify/mock
- Create temporary directories for fake Maven installations
- Clean up test files after each test
- Test both success and error paths

**Success Criteria:**
- All tests pass
- Code coverage >90%
- Tests run quickly (<1s total)
- No test pollution (cleanup works)

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 8: Unit Tests for ShimGenerator
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/generator_test.go`

**Description:** Unit tests for ShimGenerator covering shim creation, discovery, and atomic operations.

**Dependencies:** Task 3

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are writing unit tests for ShimGenerator in mvnenv-win.

**Role:** Senior Go developer with file system testing expertise.

**Task:** Create internal/shim/generator_test.go:

1. TestShimGenerator_GenerateShims_CoreCommands:
   - Create temp mvnenv directory structure
   - Create fake shim.exe binary
   - Call GenerateShims()
   - Verify: mvn.exe created
   - Verify: mvn.cmd created
   - Verify: mvnDebug.exe created
   - Verify: mvnDebug.cmd created

2. TestShimGenerator_GenerateShims_AdditionalCommands:
   - Create temp versions directory with mvnyjp.cmd
   - Call GenerateShims()
   - Verify: mvnyjp.exe created
   - Verify: mvnyjp.cmd created

3. TestShimGenerator_GenerateShims_OverwriteExisting:
   - Generate shims
   - Modify one shim file
   - Generate shims again
   - Verify: modified file overwritten

4. TestShimGenerator_GenerateShims_AtomicOperation:
   - Intercept file creation to simulate failure mid-generation
   - Call GenerateShims()
   - Verify: no partial .tmp files left

5. TestShimGenerator_DiscoverAdditionalCommands:
   - Create versions directory with multiple versions
   - Add custom commands to some versions
   - Call discoverAdditionalCommands()
   - Verify: unique commands discovered
   - Verify: core commands (mvn, mvnDebug) not included

6. TestShimGenerator_BatchShimContent:
   - Generate batch shim
   - Read content
   - Verify: calls .exe with %*
   - Verify: exits with %ERRORLEVEL%

7. TestShimGenerator_NoVersionsInstalled:
   - Create empty versions directory
   - Call GenerateShims()
   - Verify: core shims still generated
   - Verify: no error

**Restrictions:**
- Use temporary directories (t.TempDir())
- Create fake shim.exe for testing
- Verify file permissions where applicable
- Test Windows-specific behavior (.exe, .cmd)
- Clean up automatically with t.TempDir()

**Success Criteria:**
- All tests pass
- Code coverage >90%
- Both .exe and .cmd generation tested
- Atomic operations verified
- Discovery logic tested

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 9: Unit Tests for PathManager
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/path_windows_test.go`

**Description:** Unit tests for Windows PATH management including registry operations.

**Dependencies:** Task 4

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are writing unit tests for PathManager in mvnenv-win.

**Role:** Senior Go developer with Windows platform testing expertise.

**Task:** Create internal/shim/path_windows_test.go:

1. Build tag: //go:build windows

2. TestPathManager_AddShimsDirToPath:
   - Create PathManager with temp directory
   - Mock or use test registry key (careful!)
   - Call AddShimsDirToPath()
   - Verify: returns true (modified)
   - Verify: shimsDir in PATH
   - Verify: shimsDir at beginning of PATH

3. TestPathManager_AddShimsDirToPath_AlreadyPresent:
   - Add shimsDir to PATH manually
   - Call AddShimsDirToPath()
   - Verify: returns false (not modified)
   - Verify: no duplicate entries

4. TestPathManager_IsInPath_CaseInsensitive:
   - Create PATH with uppercase directory
   - Check for lowercase directory
   - Verify: returns true (case-insensitive match)

5. TestPathManager_RemoveShimsDirFromPath:
   - Add shimsDir to PATH
   - Call RemoveShimsDirFromPath()
   - Verify: returns true (modified)
   - Verify: shimsDir not in PATH

6. TestPathManager_RemoveShimsDirFromPath_NotPresent:
   - Ensure shimsDir not in PATH
   - Call RemoveShimsDirFromPath()
   - Verify: returns false (not modified)

7. TestPathManager_RemoveFromPath:
   - Create PATH with multiple entries
   - Remove one entry
   - Verify: entry removed
   - Verify: other entries preserved
   - Verify: order maintained

**Restrictions:**
- Windows-only tests (build tag)
- Be careful with registry modifications
- Use test registry keys if possible
- Restore original PATH after tests
- May require admin setup for CI

**Success Criteria:**
- All tests pass on Windows
- Case-insensitive comparison tested
- No pollution of actual PATH
- Cleanup restores original state

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 10: Integration Tests for Shim System
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `test/integration/shim_integration_test.go`

**Description:** End-to-end integration tests for the complete shim workflow.

**Dependencies:** Tasks 1-6

**Estimated Effort:** 4 hours

**_Prompt:**
```
You are writing integration tests for the shim system in mvnenv-win.

**Role:** Senior Go developer with integration testing expertise.

**Task:** Create test/integration/shim_integration_test.go:

1. TestIntegration_ShimGeneration:
   - Setup: create mvnenv directory structure
   - Install a Maven version
   - Generate shims
   - Verify: all expected shim files exist
   - Verify: shim executables are valid

2. TestIntegration_ShimExecution_GlobalVersion:
   - Setup: install Maven version, set global
   - Generate shims
   - Execute mvn shim with test arguments
   - Verify: correct Maven version executed
   - Verify: arguments passed through
   - Verify: exit code 0

3. TestIntegration_ShimExecution_LocalVersion:
   - Setup: install two Maven versions
   - Set global version 1
   - Create .maven-version file with version 2
   - Execute mvn shim
   - Verify: version 2 used (local overrides global)

4. TestIntegration_ShimExecution_ShellVersion:
   - Setup: install two versions, set global
   - Set MVNENV_MAVEN_VERSION environment variable
   - Execute mvn shim
   - Verify: shell version used (highest priority)

5. TestIntegration_ShimExecution_VersionNotInstalled:
   - Setup: set global to non-installed version
   - Execute mvn shim
   - Verify: error message mentions version not installed
   - Verify: suggests installation command
   - Verify: exit code 1

6. TestIntegration_ShimExecution_NoVersionSet:
   - Setup: no versions configured
   - Execute mvn shim
   - Verify: error message about no version set
   - Verify: suggests setting global version
   - Verify: exit code 1

7. TestIntegration_ShimIOForwarding:
   - Create test Maven script that reads stdin, writes stdout/stderr
   - Execute shim with piped input
   - Capture stdout and stderr
   - Verify: stdin forwarded correctly
   - Verify: stdout captured
   - Verify: stderr captured

8. TestIntegration_ShimExitCode:
   - Create test Maven script that exits with specific code
   - Execute shim
   - Verify: exact exit code preserved (42 -> 42)

9. TestIntegration_RehashCommand:
   - Install Maven version
   - Delete shims manually
   - Run 'mvnenv rehash'
   - Verify: shims regenerated
   - Verify: success message displayed

10. TestIntegration_AutomaticRehash:
    - Install first Maven version
    - Verify: shims created automatically
    - Install second Maven version
    - Verify: shims regenerated automatically

11. TestIntegration_DebugMode:
    - Set MVNENV_DEBUG=1
    - Execute mvn shim
    - Capture stderr
    - Verify: debug information present
    - Verify: includes version, source, path, timing

**Restrictions:**
- Use temporary directories for test installations
- Create minimal test Maven installations
- Clean up test data after each test
- Tests should be runnable offline
- Skip tests if required setup not available

**Success Criteria:**
- Full workflow tested end-to-end
- All version resolution sources tested (shell/local/global)
- I/O forwarding verified
- Exit code preservation verified
- Error cases tested
- Debug mode tested

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 11: Performance Benchmarks
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `test/benchmarks/shim_bench_test.go`

**Description:** Performance benchmarks to verify <50ms overhead requirement.

**Dependencies:** Tasks 1-2

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are writing performance benchmarks for mvnenv-win's shim system.

**Role:** Senior Go developer with performance testing expertise.

**Task:** Create test/benchmarks/shim_bench_test.go:

1. BenchmarkShimOverhead:
   - Measure total time from shim start to Maven process start
   - Use no-op Maven script for consistent timing
   - Benchmark excludes Maven execution time
   - Assert: overhead <50ms (fail if exceeded)

2. BenchmarkVersionResolution:
   - Measure time to resolve version only
   - Test with .maven-version file present
   - Assert: resolution <25ms

3. BenchmarkProcessSpawn:
   - Measure time from Maven path construction to process start
   - Use minimal Maven script
   - Assert: spawn <25ms

4. BenchmarkShimConcurrent:
   - Launch multiple shim executions concurrently
   - Measure per-shim overhead under load
   - Verify: no significant slowdown under concurrency

5. BenchmarkShimGeneration:
   - Measure time to generate all shims
   - Include both .exe and .cmd creation
   - Typical expected: <100ms for core commands

**Benchmark Output Format:**
```
BenchmarkShimOverhead-8           50    42.3 ms/op
BenchmarkVersionResolution-8     100    12.5 ms/op
BenchmarkProcessSpawn-8          100    18.7 ms/op
```

**Restrictions:**
- Use testing.B for benchmarks
- Use b.ResetTimer() after setup
- Create realistic test environment
- Fail benchmark if performance target exceeded
- Document timing requirements in comments

**Success Criteria:**
- All benchmarks run successfully
- Overhead <50ms consistently
- Resolution <25ms
- Process spawn <25ms
- Performance regression detection

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Documentation and Finalization

#### Task 12: Create Package Documentation
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/shim/doc.go`

**Description:** Package-level documentation with overview and usage examples.

**Dependencies:** All implementation tasks completed

**Estimated Effort:** 1 hour

**_Prompt:**
```
You are writing package documentation for mvnenv-win's shim system.

**Role:** Technical writer with Go documentation expertise.

**Task:** Create internal/shim/doc.go:

1. Package overview:
   - Purpose: Transparent Maven command interception
   - Key capabilities: version resolution, command forwarding, I/O pass-through
   - Architecture: Single binary serves multiple commands via filename detection

2. Usage examples:
   - Generating shims
   - Executing commands via shim
   - Configuring debug mode
   - Managing PATH

3. Code example:
```go
// Example: Generate shims for all Maven commands
generator := shim.NewShimGenerator(mvnenvRoot)
paths, err := generator.GenerateShims()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Generated %d shims\n", len(paths))

// Example: Execute Maven command with version resolution
resolver := version.NewVersionResolver(mvnenvRoot)
executor := shim.NewShimExecutor(resolver)
exitCode, err := executor.Execute("mvn", []string{"clean", "install"})
if err != nil {
    log.Fatal(err)
}
os.Exit(exitCode)

// Example: Enable debug mode
os.Setenv("MVNENV_DEBUG", "1")
// Now shims will output detailed diagnostics
```

4. Performance characteristics:
   - <50ms overhead from shim to Maven execution
   - No I/O buffering
   - Direct process forwarding

5. Debug mode documentation:
   - Set MVNENV_DEBUG=1
   - Output goes to stderr
   - Includes timing information

**Restrictions:**
- Follow godoc conventions
- Keep examples concise and runnable
- Include performance notes

**Success Criteria:**
- Package documentation renders in godoc
- Examples are accurate
- Performance characteristics documented

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 13: Create Build Script for Shim Executable
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `scripts/build-shim.bat` or `scripts/build-shim.ps1`

**Description:** Build script that compiles shim.exe and places it in bin/ directory.

**Dependencies:** Task 1

**Estimated Effort:** 1 hour

**_Prompt:**
```
You are creating a build script for mvnenv-win's shim executable.

**Role:** Build engineer with Go compilation expertise.

**Task:** Create scripts/build-shim.bat (and optionally .ps1):

1. Build script functionality:
   - Set GOOS=windows
   - Set GOARCH=amd64
   - Disable CGO: CGO_ENABLED=0
   - Compile: go build -o bin/shim.exe cmd/shim/main.go
   - Verify build succeeded
   - Report binary size

2. Optional optimizations:
   - Strip debug info: -ldflags="-s -w"
   - Enable optimizations: -trimpath

3. Build for both architectures (future):
   - amd64 (primary)
   - 386 (optional, if needed)

4. Error handling:
   - Check if Go is installed
   - Verify build succeeded
   - Print clear success/failure message

5. Output example:
```
Building shim executable...
GOOS=windows GOARCH=amd64 go build -o bin/shim.exe cmd/shim/main.go
Build successful: bin/shim.exe (2.5 MB)
```

**Restrictions:**
- Windows batch script (.bat) for portability
- Optional PowerShell version for modern environments
- No external dependencies beyond Go
- Clear error messages

**Success Criteria:**
- Script builds shim.exe successfully
- Binary placed in correct location (bin/)
- Build errors clearly reported
- Works in cmd.exe and PowerShell

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 14: Final Review and Error Message Audit
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Review:**
- All implementation files in `internal/shim/` and `cmd/shim/`

**Description:** Comprehensive review ensuring all error messages are clear, actionable, and match requirements.

**Dependencies:** All previous tasks completed

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are conducting a final review of the shim system implementation for mvnenv-win.

**Role:** Senior Go developer and technical reviewer.

**Task:** Review all implementation for completeness and quality:

1. Error Message Audit:
   - Review all error messages in internal/shim/ and cmd/shim/
   - Verify: errors include context (command, version, source)
   - Verify: errors are actionable (suggest next steps)
   - Examples:
     * "Maven version '3.9.4' is set but not installed.\nInstall it with: mvnenv install 3.9.4"
     * "No Maven version is set.\nSet a global version with: mvnenv global <version>"
     * "Maven binary not found at C:\...\.mvnenv\versions\3.9.4\bin\mvn.cmd\nVersion 3.9.4 may be corrupted. Try reinstalling with: mvnenv install 3.9.4"

2. Requirements Verification:
   - Check requirements.md against implementation
   - Verify all acceptance criteria met
   - Document any deviations

3. Performance Verification:
   - Run benchmarks: go test -bench=. ./test/benchmarks/
   - Verify: overhead <50ms
   - Verify: resolution <25ms
   - If targets not met: investigate and optimize

4. Code Quality:
   - Run go vet on shim packages
   - Run staticcheck
   - Run golangci-lint
   - Fix any issues

5. Test Coverage:
   - Run: go test -cover ./internal/shim/...
   - Verify: coverage >90%
   - Add tests for any uncovered critical paths

6. Integration Verification:
   - Verify core-version-management integration works
   - Verify automatic regeneration after install/uninstall
   - Verify manual regeneration via rehash command
   - Test debug mode (MVNENV_DEBUG=1)

7. Windows Compatibility:
   - Test in PowerShell
   - Test in cmd.exe
   - Test in Git Bash (if available)
   - Verify PATH management works

8. Documentation Review:
   - Verify all exported functions have godoc comments
   - Check package doc.go is complete
   - Ensure build scripts documented

**Success Criteria:**
- All requirements acceptance criteria met
- Error messages clear and actionable
- Performance targets achieved (<50ms)
- Code quality checks pass
- Test coverage >90%
- Windows compatibility verified
- Documentation complete

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

## Summary

**Total Tasks:** 14
**Estimated Total Effort:** 30 hours

**Task Dependencies Flow:**
```
Foundation (Tasks 1-4)
    ↓
Integration (Tasks 5-6)
    ↓
Testing (Tasks 7-11)
    ↓
Documentation (Tasks 12-14)
```

**Critical Path:**
Task 1 → Task 2 → Task 7 → Task 10 → Task 11 → Task 14

**Parallel Work Opportunities:**
- Tasks 1, 3, 4 (executor, generator, path) can be done in parallel
- Tasks 5, 6 (installer integration, rehash command) can be done in parallel after Task 3
- Tasks 7, 8, 9 (unit tests) can be done in parallel after respective implementations
- Tasks 12, 13 (documentation, build script) can be done in parallel

**Key Deliverables:**
1. Single shim executable serving all Maven commands
2. Automatic shim regeneration after install/uninstall
3. Manual regeneration via `mvnenv rehash` command
4. <50ms overhead from shim to Maven execution
5. Debug mode for troubleshooting (MVNENV_DEBUG=1)
6. Windows PATH management
7. Comprehensive test coverage (>90%)
