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
3. (WIP) Optionally fetches external metadata and caches it locally
4. Builds a normalized intermediate JSON artifact
5. Generates a Markdown recap with viewing statistics

---

## Intended Workflow

```bash
# 1. Build normalized data using cache only (no network access)
nfrecap build NetflixViewingHistory.csv

# 2. Build while fetching metadata from external APIs
nfrecap build NetflixViewingHistory.csv --fetch

# 3. Generate a recap from the built JSON
nfrecap recap NetflixViewingHistory.json --year 2025 > Netflix-2025.md
```

---

## Input

### Netflix Viewing History CSV

This tool expects the CSV file downloaded directly from Netflix.

Example:

```csv
Title,Date
"ジョン・ウィック: コンセクエンス","12/1/25"
"ジョン・ウィック: パラベラム","12/1/25"
"ジョン・ウィック：チャプター2","12/1/25"
"ジョン・ウィック","12/1/25"
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
      "date": "2025-12-01",
      "normalized": {
        "raw_title": "ジョン・ウィック: コンセクエンス",
        "work_title": "ジョン・ウィック: コンセクエンス",
        "type": "movie"
      },
      "metadata": {
        "provider": "tmdb-stub",
        "id": "stub:ジョン・ウィック: コンセクエンス",
        "title": "ジョン・ウィック: コンセクエンス",
        "genres": [
          "Unknown"
        ],
        "runtime_min": 120
      }
    },
  ]
}
```

- One file per `build` execution
- Used as the input for `recap`

---

### Recap Output (Markdown)

```md
# Netflix Recap 2025

## 視聴本数
- 合計: 4

## 月別視聴本数
- 01月: 0
- 02月: 0
- 03月: 0
- 04月: 0
- 05月: 0
- 06月: 0
- 07月: 0
- 08月: 0
- 09月: 0
- 10月: 0
- 11月: 0
- 12月: 4

## 曜日別視聴本数
- 日: 0
- 月: 4
- 火: 0
- 水: 0
- 木: 0
- 金: 0
- 土: 0

## 最長連続視聴（streak）
- 最長: 1日
```

---

## Planned / Not Yet Implemented

- Real external metadata provider implementation (e.g. TMDb)
- Genre-based estimated viewing time
- Series-level rankings
- Multi-year recaps
- Support for multiple languages (Currently supports only Japanese)

> These features are not part of the current specification and may change.

---

## License

MIT
