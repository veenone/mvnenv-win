# Product Requirements Document: mvnenv-win

## Project Overview

### Product Name
**mvnenv-win** - Maven Environment Manager for Windows

### Version
1.0.0

### Document Status
Draft

### Last Updated
October 2025

---

## 1. Executive Summary

### 1.1 Purpose
mvnenv-win is a command-line tool for Windows that manages multiple Apache Maven installations, allowing developers to easily switch between different Maven versions for various projects. It is inspired by pyenv-win but specifically designed for Maven version management.

### 1.2 Key Benefits
- **Version Flexibility**: Install and manage multiple Maven versions simultaneously
- **Project Isolation**: Set different Maven versions per project
- **Nexus Integration**: Download Maven distributions from private Nexus repositories
- **Seamless Switching**: Switch between Maven versions without system-wide configuration changes
- **Go-based Architecture**: Better performance and maintainability compared to batch scripts

### 1.3 Target Users
- Java developers working on multiple projects with different Maven requirements
- DevOps engineers managing build environments
- Teams using private Maven repositories via Nexus

---

## 2. Functional Requirements

### 2.1 Core Features

#### 2.1.1 Maven Version Management
- **Install Command**: Download and install specific Maven versions
  - Support for official Apache Maven releases
  - Support for custom Nexus repository sources
  - Silent installation option
  - Batch installation of multiple versions
  
- **Uninstall Command**: Remove installed Maven versions
  - Clean removal of Maven directories
  - Update PATH and environment variables

- **List Command**: Display available and installed versions
  - Show all available versions from configured sources
  - Filter versions by pattern
  - Indicate currently active version

#### 2.1.2 Version Selection
- **Global Version**: Set system-wide default Maven version
- **Local Version**: Set project-specific Maven version
- **Shell Version**: Temporarily override Maven version for current session

#### 2.1.3 Environment Management
- **Shim System**: Intercept Maven commands and route to correct version
- **PATH Management**: Automatically update system PATH
- **MAVEN_HOME Management**: Dynamically set MAVEN_HOME variable
- **Settings.xml Handling**: Option to maintain separate settings per version

### 2.2 Nexus Integration

#### 2.2.1 Repository Configuration
- Configure multiple Nexus repository URLs
- Authentication support (username/password, token)
- SSL/TLS certificate handling
- Repository priority ordering

#### 2.2.2 Version Discovery
- Query Nexus for available Maven versions
- Cache repository metadata locally
- Periodic update of available versions list

### 2.3 Command-Line Interface

```bash
# Core Commands
mvnenv --version              # Show mvnenv version
mvnenv --help                 # Display help information

# Version Management
mvnenv install <version>      # Install specific Maven version
mvnenv install -l            # List available versions
mvnenv uninstall <version>    # Uninstall Maven version
mvnenv versions              # List installed versions

# Version Selection
mvnenv global <version>       # Set global Maven version
mvnenv local <version>        # Set local Maven version
mvnenv shell <version>        # Set shell Maven version

# Environment Commands
mvnenv version               # Show current Maven version and origin
mvnenv which mvn            # Display path to Maven executable
mvnenv rehash               # Rebuild shim executables

# Repository Management
mvnenv repo add <name> <url>     # Add Nexus repository
mvnenv repo list                  # List configured repositories
mvnenv repo remove <name>         # Remove repository
mvnenv repo auth <name>           # Configure repository authentication

# Update Commands
mvnenv update                # Update available versions cache
mvnenv self-update           # Update mvnenv itself
```

---

## 3. Technical Requirements

### 3.1 Architecture

#### 3.1.1 Technology Stack
- **Primary Language**: Go (Golang) 1.21+
- **Build System**: Go modules
- **Configuration Format**: YAML/TOML for settings
- **Version File Format**: Plain text (.maven-version files)

#### 3.1.2 Components

```
mvnenv-win/
├── cmd/                    # Command-line interface
│   ├── mvnenv/            # Main executable
│   └── shim/              # Shim executable
├── internal/              # Internal packages
│   ├── config/           # Configuration management
│   ├── download/         # Download manager
│   ├── environment/      # Environment variable handling
│   ├── nexus/           # Nexus repository client
│   ├── shim/            # Shim generation
│   └── version/         # Version management
├── pkg/                   # Public packages
│   └── maven/           # Maven-specific utilities
└── test/                  # Test suites
```

### 3.2 Installation Structure

```
%USERPROFILE%/
└── .mvnenv/
    ├── bin/               # mvnenv executables
    ├── shims/            # Shim executables
    ├── versions/         # Installed Maven versions
    │   ├── 3.6.3/
    │   ├── 3.8.6/
    │   └── 3.9.4/
    ├── cache/            # Downloaded archives
    ├── config/           # Configuration files
    │   ├── config.yaml
    │   └── repositories.yaml
    └── logs/             # Operation logs
```

### 3.3 Configuration Schema

#### 3.3.1 Main Configuration (config.yaml)
```yaml
version: 1.0
global_version: "3.9.4"
auto_rehash: true
download_timeout: 300
proxy:
  http: ""
  https: ""
  no_proxy: ""
logging:
  level: "info"
  file: "mvnenv.log"
```

#### 3.3.2 Repository Configuration (repositories.yaml)
```yaml
repositories:
  - name: "apache"
    url: "https://archive.apache.org/dist/maven/maven-3/"
    priority: 1
    enabled: true
  - name: "company-nexus"
    url: "https://nexus.company.com/repository/maven-releases/"
    priority: 2
    auth:
      type: "basic"
      username: "${NEXUS_USER}"
      password: "${NEXUS_PASS}"
    enabled: true
```

---

## 4. Non-Functional Requirements

### 4.1 Performance
- Version switching: < 100ms
- Maven installation: Limited by network speed
- Shim execution overhead: < 50ms
- Memory footprint: < 50MB

### 4.2 Compatibility
- **Windows Versions**: Windows 10+ (64-bit)
- **PowerShell**: 5.1+
- **Command Prompt**: Full support
- **Terminal Emulators**: Windows Terminal, ConEmu, cmder

### 4.3 Security
- Checksum verification for downloaded files
- Secure credential storage using Windows Credential Manager
- SSL/TLS verification for HTTPS connections
- No elevation required for user-level installation

### 4.4 Reliability
- Atomic version installations (rollback on failure)
- Graceful handling of network interruptions
- Version integrity checks
- Backup of previous configurations

---

## 5. Installation Methods

### 5.1 Primary Installation Methods

#### 5.1.1 PowerShell Script (Recommended)
```powershell
Invoke-WebRequest -UseBasicParsing -Uri "https://raw.githubusercontent.com/veenone/mvnenv-win/master/install-mvnenv-win.ps1" -OutFile "./install-mvnenv-win.ps1"
./install-mvnenv-win.ps1
```

#### 5.1.2 Go Install
```bash
go install github.com/veenone/mvnenv-win/cmd/mvnenv@latest
```

#### 5.1.3 Manual Installation
1. Download latest release from GitHub
2. Extract to `%USERPROFILE%\.mvnenv`
3. Add `%USERPROFILE%\.mvnenv\bin` and `%USERPROFILE%\.mvnenv\shims` to PATH

### 5.2 Package Managers

#### 5.2.1 Chocolatey
```bash
choco install mvnenv-win
```

#### 5.2.2 Scoop
```bash
scoop install mvnenv-win
```

---

## 6. Development Roadmap

### Phase 1: Core Functionality (v1.0.0)
- [ ] Basic version management (install, uninstall, list)
- [ ] Global and local version selection
- [ ] Shim system implementation
- [ ] Apache Maven repository support
- [ ] PowerShell installation script

### Phase 2: Nexus Integration (v1.1.0)
- [ ] Nexus repository configuration
- [ ] Authentication mechanisms
- [ ] Custom repository support
- [ ] Version discovery from Nexus

### Phase 3: Enhanced Features (v1.2.0)
- [ ] Settings.xml per-version management
- [ ] JDK version compatibility checking
- [ ] Plugin for popular IDEs (IntelliJ, VS Code)
- [ ] Import/export configuration

### Phase 4: Advanced Features (v2.0.0)
- [ ] Cross-platform support (Linux, macOS)
- [ ] Container integration (Docker, WSL2)
- [ ] Team configuration sharing
- [ ] Automated version selection based on pom.xml

---

## 7. Testing Strategy

### 7.1 Unit Tests
- Command parsing and validation
- Version comparison logic
- Configuration management
- Shim generation

### 7.2 Integration Tests
- Maven installation process
- Version switching
- Nexus repository interaction
- Environment variable management

### 7.3 End-to-End Tests
- Complete installation workflow
- Multi-version project scenarios
- Repository failover
- Upgrade/downgrade scenarios

---

## 8. Documentation Requirements

### 8.1 User Documentation
- Installation guide
- Command reference
- Configuration guide
- Troubleshooting guide
- Migration from manual Maven management

### 8.2 Developer Documentation
- Architecture overview
- API documentation
- Contributing guidelines
- Plugin development guide

---

## 9. Success Metrics

### 9.1 Adoption Metrics
- Number of GitHub stars
- Download count
- Active installations

### 9.2 Quality Metrics
- Installation success rate > 95%
- Version switch success rate > 99%
- User-reported bug count < 10 per release
- Performance benchmarks met

### 9.3 User Satisfaction
- GitHub issue resolution time < 7 days
- Documentation completeness score > 90%
- Community contribution rate

---

## 10. Risks and Mitigation

### 10.1 Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Maven download server unavailability | High | Medium | Cache distributions locally, support mirrors |
| Windows Defender false positives | Medium | Low | Code signing certificate |
| PATH length limitations | Medium | Low | Short paths, junction points |
| Shim conflicts with antivirus | Medium | Medium | Whitelist documentation |

### 10.2 Adoption Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| User resistance to new tool | High | Medium | Clear migration guide, automation scripts |
| Corporate security policies | High | Low | Enterprise features, compliance docs |
| Competition from existing tools | Medium | Medium | Unique features, better performance |

---

## 11. Dependencies

### 11.1 External Dependencies
- Apache Maven distribution files
- Go standard library
- Selected Go packages:
  - `github.com/spf13/cobra` (CLI framework)
  - `github.com/spf13/viper` (Configuration)
  - `github.com/go-resty/resty/v2` (HTTP client)
  - `gopkg.in/yaml.v3` (YAML parsing)

### 11.2 System Dependencies
- Windows registry access
- File system permissions
- Network connectivity
- Windows Credential Manager (optional)

---

## 12. Appendices

### A. Comparison with Existing Tools

| Feature | mvnenv-win | mvnw | SDKMAN (Windows) |
|---------|------------|------|------------------|
| Windows Native | ✓ | ✓ | ✗ |
| No WSL/Cygwin | ✓ | ✓ | ✗ |
| Multiple Versions | ✓ | ✗ | ✓ |
| Project-specific | ✓ | ✓ | ✓ |
| Nexus Support | ✓ | ✗ | ✗ |
| Go-based | ✓ | ✗ | ✗ |

### B. Migration Path from pyenv-win

Users familiar with pyenv-win will find mvnenv-win commands intuitive:

| pyenv-win | mvnenv-win |
|-----------|------------|
| `pyenv install 3.9.0` | `mvnenv install 3.9.4` |
| `pyenv global 3.9.0` | `mvnenv global 3.9.4` |
| `pyenv local 3.8.0` | `mvnenv local 3.8.6` |
| `pyenv versions` | `mvnenv versions` |

---

## Document Control

- **Author**: Development Team
- **Reviewers**: Technical Lead, Product Owner
- **Approval**: Pending
- **Next Review**: Q1 2026