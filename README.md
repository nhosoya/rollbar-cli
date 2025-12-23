# rollbar-cli

A lightweight CLI tool to query Rollbar errors.

## Installation

```bash
go install github.com/nhosoya/rollbar-cli@latest
```

Or build from source:

```bash
git clone https://github.com/nhosoya/rollbar-cli.git
cd rollbar-cli
go build -o rollbar-cli
```

## Configuration

Set your Rollbar read access token as an environment variable:

```bash
export ROLLBAR_READ_TOKEN="your_token_here"
```

## Usage

### List items

```bash
# List active error items (default: 10 items)
rollbar-cli items

# List with filters
rollbar-cli items -n 5 -s active -l error -e production
```

Flags:
- `-n, --limit`: Number of items (default: 10)
- `-s, --status`: Filter by status: `active`, `resolved`, `muted` (default: "active")
- `-l, --level`: Filter by level: `error`, `warning`, `critical`
- `-e, --env`: Filter by environment

### Show item details

```bash
rollbar-cli item <item_id>
```

### List occurrences

```bash
# List occurrences for an item
rollbar-cli occurrences <item_id>

# Or use the alias
rollbar-cli occ <item_id> -n 5
```

Flags:
- `-n, --limit`: Number of occurrences (default: 10)

### Show occurrence details

```bash
# Show occurrence with essential fields
rollbar-cli occurrence <occurrence_id>

# Or use the alias with full output
rollbar-cli o <occurrence_id> -f
```

Flags:
- `-f, --full`: Show full occurrence data (raw API response)

## Output Format

All commands output JSON. Example:

```bash
$ rollbar-cli items -n 1
[
  {
    "id": "1722412686",
    "counter": 21195,
    "title": "SomeError: something went wrong",
    "level": "error",
    "status": "active",
    "environment": "production",
    "total_occurrences": 100,
    "last_occurrence": "2025-12-23 12:36:08",
    "first_occurrence": "2025-12-03 15:15:28"
  }
]
```

## License

MIT
