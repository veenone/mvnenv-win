# Nexus Repository Integration

mvnenv-win supports downloading Maven distributions from private Nexus Repository Manager instances in addition to the public Apache Maven archives.

## Features

- Download Maven distributions from private Nexus repositories
- Support for basic authentication (username/password)
- Custom CA certificate support for self-signed certificates
- Option to disable SSL certificate verification
- Automatic fallback to Apache archive if Nexus fails
- Version discovery from Nexus maven-metadata.xml

## Configuration

### Basic Setup

Create or edit your configuration file at `%USERPROFILE%\.mvnenv\config\config.yaml`:

```yaml
version: "1.0"
global_version: "3.9.4"
auto_rehash: true

nexus:
  enabled: true
  base_url: "https://nexus.example.com/repository/maven-public"
  username: "your-username"
  password: "your-password"
```

### With Self-Signed Certificate

If your Nexus server uses a self-signed certificate:

```yaml
nexus:
  enabled: true
  base_url: "https://nexus.internal.company.com/repository/maven-central"
  username: "myuser"
  password: "mypassword"
  tls:
    insecure_skip_verify: true
```

### With Custom CA Certificate

For enterprises with internal CA certificates:

```yaml
nexus:
  enabled: true
  base_url: "https://nexus.internal.company.com/repository/maven-central"
  username: "myuser"
  password: "mypassword"
  tls:
    ca_file: "C:\\company\\ca\\root-ca.pem"
```

### Without Authentication

For public Nexus repositories that don't require authentication:

```yaml
nexus:
  enabled: true
  base_url: "https://nexus.public.company.com/repository/maven-central"
```

## Usage

Once configured, mvnenv will automatically use Nexus when available:

```bash
# Update version cache (fetches from Nexus + Apache)
mvnenv update

# List available versions (includes Nexus versions)
mvnenv install -l

# Install from Nexus (tries Nexus first, falls back to Apache)
mvnenv install 3.9.4

# Check latest version (checks Nexus + Apache)
mvnenv latest --remote
```

## Repository Behavior

mvnenv uses a **failover strategy** when both Nexus and Apache are configured:

1. **Version Discovery**: Combines versions from both Nexus and Apache archives
2. **Download Priority**: Attempts Nexus first, falls back to Apache if Nexus fails
3. **Automatic Fallback**: If Nexus is unavailable, mvnenv continues with Apache archives

This ensures maximum availability even if your Nexus server is temporarily down.

## Nexus Repository Requirements

Your Nexus repository must:

1. Be a **Maven 2 repository** (proxy or hosted)
2. Contain Maven distributions at the standard path:
   ```
   org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip
   ```
3. Provide `maven-metadata.xml` at:
   ```
   org/apache/maven/apache-maven/maven-metadata.xml
   ```

### Recommended Nexus Setup

For best results, configure Nexus as a **proxy repository** pointing to:
- `https://repo.maven.apache.org/maven2/`

This allows Nexus to cache Maven distributions from Apache while providing your authentication and SSL requirements.

## Troubleshooting

### Error: "failed to fetch metadata"

- Verify the `base_url` is correct
- Check that the repository is a Maven 2 repository
- Ensure the repository contains Apache Maven artifacts

### Error: "authentication failed"

- Verify username and password are correct
- Check that your Nexus user has read access to the repository
- Ensure the credentials don't contain special YAML characters (quote them if needed)

### Error: "certificate verification failed"

- Use `insecure_skip_verify: true` for self-signed certificates
- Or add your CA certificate with `ca_file` option
- Ensure the CA certificate file path is correct and accessible

### Nexus not being used

- Check that `enabled: true` is set
- Verify the configuration file is at the correct location
- Run `mvnenv update` to refresh the version cache
- Check for warning messages during install/update operations

## Security Considerations

### Credential Storage

Currently, credentials are stored in plain text in the configuration file. For production use:

1. Store the config file in a secure location
2. Set appropriate file permissions (read-only for your user)
3. Consider using environment variables (future feature)

### SSL/TLS

- **Production**: Use `ca_file` to specify your company's CA certificate
- **Development/Testing**: Use `insecure_skip_verify: true` only in trusted networks
- **Never** disable certificate verification when accessing Nexus over the internet

## Example Configuration File

See `config.example.yaml` in the repository root for a complete configuration example.

## Future Enhancements

Planned features for Nexus integration:

- Windows Credential Manager integration for secure credential storage
- Multiple repository support with priority ordering
- Token-based authentication (bearer tokens)
- Repository health checking
- Per-repository timeout configuration
