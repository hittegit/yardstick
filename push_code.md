# Controlled Promotion Command (Yardstick)

## Purpose

Safely move changes through the Yardstick delivery pipeline:

```text
Issue -> Feature Branch -> Pull Request -> main -> Release Tag
```

With:

- Issue-first planning and acceptance criteria
- Issue-numbered branch naming
- Pull-request based integration (no direct pushes to `main`)
- Explicit confirmation before merge and release tagging
- Repository-specific validation gates for Go and config/docs quality

## Global Safety Rules (Non-Negotiable)

1. No direct pushes to `main`
2. Every non-trivial change starts with a GitHub Issue
3. Branch names must include the issue number
4. All integrations happen via Pull Request
5. Sync remotes before branch/release logic
6. Use remote refs as source of truth (`origin/main`)
7. Require explicit confirmation before merge or release tag push

## Repository Validation Requirements

Run these before marking a PR ready for review:

```bash
go test ./...
go vet ./...
markdownlint '*.md'
yamllint .golangci.yml .goreleaser.yaml .github/workflows/*.yml .markdownlint.yaml
```

Hard abort on any failure.

## Step 0: Pre-Flight Checks (All Modes)

### 1. Verify repository context

```bash
git remote -v
git rev-parse --show-toplevel
```

Expected repo: `hittegit/yardstick`.

### 2. Sync remotes

```bash
git fetch origin --prune --tags
```

### 3. Verify GitHub auth

```bash
gh auth status
```

If auth fails, abort operations requiring issue/PR/release API calls.

## Command Modes

### Mode A: `/push_code start`

Create or use an issue, create branch, and open a draft PR.

#### Mode A Steps

1. Confirm issue exists (or create one first):

   ```bash
   gh issue create --title "<title>" --body-file <file>
   ```

2. Create and switch to issue branch:

   ```bash
   git checkout -b feat/<issue-number>-<short-slug>
   ```

3. Push branch:

   ```bash
   git push -u origin feat/<issue-number>-<short-slug>
   ```

4. Open draft PR linked to issue:

   ```bash
   gh pr create \
     --base main \
     --head feat/<issue-number>-<short-slug> \
     --draft \
     --title "<type>: <summary>" \
     --body "Closes #<issue-number>"
   ```

### Mode B: `/push_code update`

Commit and push current branch updates to the open PR.

#### Mode B Steps

1. Ensure not on protected branch:

   ```bash
   BRANCH=$(git rev-parse --abbrev-ref HEAD)
   ```

   Abort if branch is `main`.

2. Show working changes:

   ```bash
   git status --short
   git diff --stat
   ```

3. Generate proposed commit message (Conventional Commit style):

   ```text
   type(scope): short summary
   ```

4. Confirmation gate for commit message:

   ```text
   Proposed commit shown to user
   [Y]es / [n]o / [e]dit
   ```

5. If approved:

   ```bash
   git add -A
   git commit -m "<approved-message>"
   git push
   ```

### Mode C: `/push_code ready`

Validate branch, mark PR ready for review, and request review.

#### Mode C Steps

1. Run validation suite:

   ```bash
   go test ./...
   go vet ./...
   markdownlint '*.md'
   yamllint .golangci.yml .goreleaser.yaml .github/workflows/*.yml .markdownlint.yaml
   ```

2. If all pass, convert PR from draft:

   ```bash
   gh pr ready
   ```

3. Optionally request reviewers:

   ```bash
   gh pr edit --add-reviewer <user>
   ```

### Mode D: `/push_code merge`

Merge PR into `main` only after explicit confirmation.

#### Mode D Confirmation Gate

```text
This will merge into main.
Type YES to continue.
```

If response is not exactly `YES`, abort.

#### Mode E Steps

1. Confirm PR checks are green:

   ```bash
   gh pr checks
   ```

2. Merge (squash recommended):

   ```bash
   gh pr merge --squash --delete-branch
   ```

### Mode E: `/push_code release`

Cut and push a new semver tag from `main` after merge.

#### Confirmation gate

```text
This will create and push a production release tag from main.
Type YES to continue.
```

If response is not exactly `YES`, abort.

#### Steps

1. Update local `main`:

   ```bash
   git checkout main
   git pull --ff-only
   ```

2. Select next version (`vX.Y.Z`) based on merged scope.

3. Create and push annotated tag:

   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin vX.Y.Z
   ```

4. Verify release workflow starts (`.github/workflows/release.yml`).

## Abort Conditions

Immediately abort if any of the following is true:

- Wrong repository detected
- GitHub auth missing for API operations
- Validation commands fail
- User does not pass confirmation gate for merge/release
- Attempted operation targets `main` without PR flow

## PR Description Template

```md
## Summary
- <change 1>
- <change 2>

## Validation
- go test ./...
- go vet ./...
- markdownlint '*.md'
- yamllint .golangci.yml .goreleaser.yaml .github/workflows/*.yml .markdownlint.yaml

## Tracking
Closes #<issue-number>
```

## Notes Specific to Yardstick

- Yardstick is a read-only scanner; changes should not introduce write behavior in checks.
- Preserve JSON contract stability unless a versioned release note explicitly states changes.
- Keep CI-focused behavior deterministic and local-file based.
