# nfrecap

`nfrecap` is a CLI tool that generates **recap reports** from Netflix viewing history CSV files.

- Written in Go
- Distributed as a single binary
- Offline-first by default
- External API access is **explicitly opt-in**

---

## What This Tool Does

1. Reads Netflix viewing history CSV (`Title`, `Date`)
2. Normalizes titles (movie / TV, work title, season, episode)
3. Optionally fetches external metadata and caches it locally
4. Builds a normalized intermediate JSON artifact
5. Generates a Markdown recap with viewing statistics

---

## Intended Workflow

```bash
# 1. Build normalized data using cache only (no network access)
nfrecap build --in NetflixViewingHistory.csv --out NetflixViewingHistory.json

# 1. Build while fetching metadata from external APIs
TMDB_BEARER_TOKEN=<your token> nfrecap build --in NetflixViewingHistory.csv --out NetflixViewingHistory.json --fetch

# 2. Generate a recap from the built JSON
nfrecap recap --in NetflixViewingHistory.json --year 2025 --out Netflix-2025.md
```

---

## Input

### Netflix Viewing History CSV

This tool expects the CSV file downloaded directly from Netflix.

Example:

```csv
Title,Date
"駒田蒸留所へようこそ","6/7/25"
```

- Required columns: `Title`, `Date`
- Date format: `M/D/YY` (e.g. `12/13/25`)

---

## Commands

### `nfrecap build`

Builds a **normalized JSON dataset** from Netflix viewing history.

```bash
nfrecap build --in NetflixViewingHistory.csv --out NetflixViewingHistory.json
```

#### Options

| Option        | Description                                              |
| ------------- | -------------------------------------------------------- |
| `--fetch`     | Fetch metadata from external APIs and update the cache   |
| `--cache-dir` | Metadata cache directory (default: OS cache directory)   |
| `--out`       | Output JSON file (default: `NetflixViewingHistory.json`) |
| `--verbose`   | Enable verbose logging                                   |

#### Behavior

- **Default (without `--fetch`)**
  - No network access
  - Uses cached metadata only
  - Items missing metadata are marked as unresolved

- **With `--fetch`**
  - Fetches metadata from external APIs
  - Saves results to the local cache
  - Future runs can reuse the cache without `--fetch`

---

### `nfrecap recap`

Generates a **Markdown recap** from a previously built JSON file.

```bash
nfrecap recap --in NetflixViewingHistory.json --year 2025 --out Netflix-2025.md
```

#### Currently Generated Statistics

- Total number of views
- Views by month
- Views by weekday
- Longest viewing streak
  - Number of consecutive days with at least one view
  - Start and end dates of the streak

---

## Output Formats

### Build Output (JSON)

```json
{
  "generated_at": "2025-12-22T14:39:53+09:00",
  "items": [
    {
      "date": "2025-06-07",
      "normalized": {
        "raw_title": "駒田蒸留所へようこそ",
        "work_title": "駒田蒸留所へようこそ",
        "type": "movie"
      },
      "metadata": {
        "provider": "tmdb",
        "id": "movie:1119211",
        "title": "Komada – A Whisky Family",
        "genres": [
          "Animation",
          "Drama",
          "Family"
        ],
        "runtime_min": 91
      }
    },
  ]
}
```

- One file per `build` execution
- Used as the input for `recap`

---

### Recap Output (Markdown)

See [testdata/recap-sample.md](testdata/recap-sample.md)

---

## Planned / Not Yet Implemented

- Multi-year recaps
- Support for multiple languages (Currently supports only Japanese)

> These features are not part of the current specification and may change.

---

## Disclaimer

This nfrecap uses TMDB and the TMDB APIs but is not endorsed, certified, or otherwise approved by TMDB.

---

## License

MIT
