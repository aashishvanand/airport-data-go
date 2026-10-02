---
name: update-airport-data
description: Sync data/airports.json from airport-data-js into airport-data-go, regenerate derived data files, check dependencies, bump the patch version and release. Use when the user says the JS library has new airport data, asks to update/sync airports, or asks for a data release.
---

# Update airport data (airport-data-go)

`airport-data-js` is the source of truth for the dataset. This repo ships a copy of it.
Work from the repo root (`/Users/aashishvanand/Code/airport-data-go`).

## 1. Get the new data

```bash
git checkout main && git pull --ff-only
git -C ../airport-data-js fetch origin
git -C ../airport-data-js show origin/main:data/airports.json > data/airports.json
git diff --stat data/airports.json   # nothing changed? stop, there is nothing to release
```

## 2. Scan field types before touching code

Upstream sometimes changes value types (for example, `""` became `null` for missing
`elevation_ft` / `runway_length` in JS 4.0.0). Compare old against new for every field:

```bash
git show HEAD:data/airports.json > /tmp/old_airports.json
python3 - <<'EOF'
import json
from collections import Counter
o = json.load(open('/tmp/old_airports.json')); n = json.load(open('data/airports.json'))
print('count', len(o), '->', len(n))
for f in n[0]:
    co = Counter(type(a.get(f)).__name__ for a in o); cn = Counter(type(a.get(f)).__name__ for a in n)
    if co.keys() != cn.keys():
        print(f, dict(co), '->', dict(cn))
old = {(a['iata'], a['icao']) for a in o}
print('added', [(a['iata'], a['icao'], a['airport']) for a in n if (a['iata'], a['icao']) not in old])
EOF
```

If a field gains a new type (`NoneType`, `str` in a numeric field, and so on), check that the decoder below handles it
and add a test for it. Also read the top of `../airport-data-js/CHANGELOG.md` for data fixes and breaking type changes.

Decoder: `rawAirport` uses `json.RawMessage`, together with `rawToFloat64` / `rawToInt` in `airport_data.go`.
These handle numbers, numeric strings, `""`, and `null`. A null in a string field would still break `Airport`.

## 3. Regenerate derived files

Nothing to do here. `data/airports.json` is embedded with `go:embed`.

## 4. Check for dependency updates

There are no module dependencies, so just run `go mod verify`. Check `.github/workflows/*.yml` action versions
(they are pinned to SHAs) and merge any open dependabot PRs first (`gh pr list`).

## 5. Bump the version (patch for data-only updates)

Write the new version to the `VERSION` file, for example `echo "1.0.3" > VERSION`.

## 6. Verify locally (same checks as CI)

```bash
go mod verify
go vet ./...
go test -race ./...
```

If a test fails, confirm the failure comes from a real data change (new airport count, corrected
coordinates, and so on) before you update the expected value.

## 7. Commit and push main, then wait for CI

Stage explicit paths only. Never use `git add .`, because there are untracked local files.

```bash
git add VERSION data/airports.json
git commit -m "feat: sync airport data with airport-data-js and bump version to X.Y.Z"
git push origin main
gh run list --branch main --limit 1        # then: gh run watch <id> --exit-status
```

## 8. Release (only after CI on main is green; publishing is irreversible)

Merging into `release` makes the workflow create the `vX.Y.Z` tag and GitHub Release, and index the module on proxy.golang.org:
```bash
git checkout release && git pull --ff-only && git merge --ff-only main && git push origin release
git checkout main
```
Check the result with `curl -s https://proxy.golang.org/github.com/aashishvanand/airport-data-go/@v/list`.
