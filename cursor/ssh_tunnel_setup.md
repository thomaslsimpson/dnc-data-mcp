# SSH Tunnel Setup Documentation

## Connection Details
- SSH Host: blackhole.dnc.io
- SSH Port: 22
- SSH User: (from config)
- SSH Key Path: (from config)
- Local Listener: 127.0.0.1:5433
- Remote Database: (from config)

## Running the MCP Service
IMPORTANT: Always run the service from the project root directory (`/Users/tsimpson/github/dnc-data-mcp`):
```bash
# Correct way to run the service
cd /Users/tsimpson/github/dnc-data-mcp
GIN_MODE=debug go run main.go

# DO NOT run from parent directory
cd ..  # This will fail
go run main.go
```

## Database Configuration
IMPORTANT: When using the SSH tunnel:
1. The tunnel uses the original database server/port from config
2. The database connection should use localhost:5433
3. Create a copy of the config for the database connection to avoid modifying the tunnel config

## Common Issues and Solutions

### 1. Connection Refused Error
If you see: `Remote dial error: ssh: rejected: connect failed (Connection refused)`
- The remote database is not running on the expected host/port
- The SSH tunnel is trying to connect to the wrong address
- The remote database is not accepting connections

### 2. SSH Key Issues
If you see: `unable to read private key` or `unable to parse private key`
- Ensure the SSH key exists at the specified path
- Check key permissions (should be 600)
- Verify the key is in the correct format

### 3. Database Connection Issues
If you see: `error connecting to the database: read tcp 127.0.0.1:XXXXX->127.0.0.1:5433: read: connection reset by peer`
- The SSH tunnel is not properly forwarding to the remote database
- The remote database might be rejecting the connection
- The local port might be in use

## Configuration
Database and SSH configuration should be stored in `~/.ssh/dnc_db_info`. This file contains:
- SSH connection details
- Database connection details
- Required credentials

## Verification Steps
1. SSH connection works: `ssh -i <key_path> <user>@blackhole.dnc.io`
2. Remote database is accessible: `ssh -i <key_path> blackhole.dnc.io "nc -zv <db_host> <db_port>"`
3. Local port is available: `lsof -i :5433`
4. SSH tunnel is forwarding: `netstat -an | grep 5433`

## Important Notes
- The SSH tunnel forwards local port 5433 to the remote database
- All database connections should be made to localhost:5433
- The SSH key must have the correct permissions (600)
- Never commit sensitive information to version control
- Always use configuration files for credentials and connection details
- Always run the service from the project root directory
- Never modify the original config object when setting up the database connection