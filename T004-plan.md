# T004: Evaluate removing binary from git history

## Objective
Assess the repository size and impact of the committed binary on history, and determine if using `git-filter-repo` (or similar) is necessary and safe.

## Findings

### Repository Analysis
- Current repository size: 4.7 MB (.git directory)
- Binary file: tools/coverage-check/coverage-check
- Binary size: 3,059,042 bytes (~3 MB)
- Binary added: d8a07ac (Tue Apr 22 17:05:33 2025) - "feat: create Go coverage checker tool"
- Binary removed: f0eb7a6 (Wed Apr 23 09:15:36 2025) - "chore: unstage binary executable"
- Only 2 commits involved with the binary
- Repository has 3 branches (master, feature/convert-to-proper-package, feature/output-file-option)

### Impact Assessment
- The binary represents a significant portion of the repository size (~64% of the .git directory size)
- The binary was only tracked for a brief period (less than 1 day)
- The binary is already properly excluded from tracking going forward via .gitignore
- The binary is properly built in CI as needed
- Removing the binary would reduce clone and fetch sizes considerably

### Risks of Rewriting History
- All branches would need to be rebased onto the new history
- All collaborators would need to re-clone or update their local repos
- PRs based on old history might break
- Repository link references might break
- Potential loss of historical context

## Recommendation

Given that:
1. The binary represents a significant portion of the repository size (~3 MB out of 4.7 MB)
2. The repository is still relatively small, and the binary was only present briefly in history
3. The risk of rewriting history is moderate, affecting all branches and collaborators
4. The project appears to be in early stages with manageable collaboration complexity

**Recommendation**: Create a separate ticket to implement the history rewrite using `git-filter-repo` to remove the binary. The size reduction justifies the effort, especially if project development is accelerating and additional contributors are expected.

## Implementation Plan (for the follow-up ticket)
1. Notify all contributors about the planned history rewrite
2. Create a backup of the repository
3. Use `git-filter-repo` to remove the binary file from history
4. Force-push the rewritten history to remote
5. Have all contributors re-clone the repository or perform a careful update
6. Update any open PRs as needed
