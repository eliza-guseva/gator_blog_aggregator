# Gator - RSS Feed Aggregator

Gator is a command-line RSS feed aggregator built in Go. It allows users to manage RSS feeds, follow/unfollow feeds, and browse posts from their subscribed feeds in the terminal.

## Features

- **User Management**: Register and login users with persistent configuration
- **Feed Management**: Add, list, and manage RSS feeds
- **Feed Following**: Follow and unfollow specific RSS feeds
- **RSS Aggregation**: Automatically fetch and store posts from RSS feeds
- **Post Browsing**: View posts from followed feeds in chronological order
- **PostgreSQL Integration**: Persistent storage using PostgreSQL database
- **CLI Interface**: Simple command-line interface for all operations

## Architecture

The project follows a clean architecture pattern with the following components:

- **CLI Commands** (`internal/cmd/`): Command handlers and business logic
- **Database Layer** (`internal/database/`): Generated database queries using sqlc
- **RSS Parser** (`rss/`): RSS feed fetching and parsing
- **Configuration** (`internal/config/`): User configuration management
- **Database Migrations** (`sql/schema/`): Database schema versioning

## Database Schema

The application uses PostgreSQL with the following main entities:

- **Users**: User accounts with unique names
- **Feeds**: RSS feed URLs with metadata
- **Feed Follows**: Many-to-many relationship between users and feeds
- **Posts**: Individual RSS feed entries

## Prerequisites

- Go 1.24.5 or later
- PostgreSQL database
- sqlc (for database code generation)
- goose (for database migrations)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd gator
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database and configure connection string

4. Run database migrations:
```bash
goose -dir sql/schema postgres "your-connection-string" up
```

5. Generate database code:
```bash
sqlc generate
```

6. Build the application:
```bash
go build -o gator
```

## Configuration

Create a configuration file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost/gator?sslmode=disable",
  "current_user_name": ""
}
```

## Usage

### User Management

**Register a new user:**
```bash
./gator register <username>
```

**Login as existing user:**
```bash
./gator login <username>
```

**List all users:**
```bash
./gator users
```

**Reset all users (development):**
```bash
./gator reset
```

### Feed Management

**Add a new RSS feed:**
```bash
./gator addfeed <feed-name> <feed-url>
```

**List all feeds:**
```bash
./gator feeds
```

### Feed Following

**Follow a feed:**
```bash
./gator follow <feed-url>
```

**List feeds you're following:**
```bash
./gator following
```

**Unfollow a feed:**
```bash
./gator unfollow <feed-url>
```

### RSS Aggregation

**Start the aggregation service:**
```bash
./gator agg <duration>
```

Example: `./gator agg 30s` fetches feeds every 30 seconds

### Browse Posts

**Browse posts from followed feeds:**
```bash
./gator browse [limit]
```

Example: `./gator browse 10` shows the 10 most recent posts

## Commands Reference

| Command | Description | Authentication Required |
|---------|-------------|------------------------|
| `register <username>` | Register a new user | No |
| `login <username>` | Login as existing user | No |
| `reset` | Delete all users (dev only) | No |
| `users` | List all users | No |
| `agg <duration>` | Start RSS aggregation service | No |
| `addfeed <name> <url>` | Add and follow a new RSS feed | Yes |
| `feeds` | List all RSS feeds | No |
| `follow <url>` | Follow an existing RSS feed | Yes |
| `following` | List feeds you're following | Yes |
| `unfollow <url>` | Unfollow a RSS feed | Yes |
| `browse [limit]` | Browse posts from followed feeds | Yes |

## Development

### Database Migrations

Add new migrations in `sql/schema/` following the naming convention:
```
XXX_description.sql
```

### Adding New Queries

1. Add SQL queries in `sql/queries/`
2. Run `sqlc generate` to generate Go code
3. Use the generated functions in your handlers

### Project Structure

```
gator/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── sqlc.yaml              # sqlc configuration
├── internal/
│   ├── cmd/               # CLI command handlers
│   ├── config/            # Configuration management
│   └── database/          # Generated database code
├── rss/                   # RSS parsing functionality
└── sql/
    ├── queries/           # SQL queries for sqlc
    └── schema/            # Database migrations
```

## Dependencies

- **github.com/lib/pq**: PostgreSQL driver for Go
- **github.com/google/uuid**: UUID generation and parsing
- **sqlc**: SQL code generation
- **goose**: Database migration tool

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of a learning exercise and is provided as-is for educational purposes.

## Troubleshooting

### Common Issues

1. **Database connection errors**: Verify PostgreSQL is running and connection string is correct
2. **Migration errors**: Ensure goose is installed and migrations are run in order
3. **Permission errors**: Check file permissions for config file creation
4. **RSS parsing errors**: Some feeds may have non-standard date formats

### Debug Mode

For debugging, you can add verbose logging or examine the database directly:

```sql
-- View all users
SELECT * FROM users;

-- View all feeds with user info
SELECT f.name, f.url, u.name as owner FROM feeds f JOIN users u ON f.user_id = u.id;

-- View recent posts
SELECT p.title, p.published_at, f.name as feed_name 
FROM posts p JOIN feeds f ON p.feed_id = f.id 
ORDER BY p.published_at DESC LIMIT 10;
```
