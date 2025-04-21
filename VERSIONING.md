# Versioning Strategy

This document outlines the versioning strategy for the Handoff project, which follows Semantic Versioning (SemVer) principles.

## Initial Version

The initial release version for Handoff will be **v0.1.0**.

### Rationale

The v0.1.0 starting point was chosen for the following reasons:

1. **Recently Refactored API**: The codebase has recently undergone significant refactoring to create a proper library/CLI separation. While the implementation is functional, the API may benefit from real-world usage feedback before committing to stability.

2. **Testing Maturity**: Although test coverage is good (85%+), additional real-world usage will help identify edge cases and potential improvements.

3. **API Evolution**: Starting with a v0.x version allows us to make necessary adjustments to the API based on user feedback without breaking semantic versioning principles.

4. **Go Convention**: Many Go libraries begin with v0.x versions until their APIs stabilize.

## Semantic Versioning

Handoff follows [Semantic Versioning 2.0.0](https://semver.org/) with the following version structure:

```
vMAJOR.MINOR.PATCH
```

- **MAJOR**: Incremented for incompatible API changes
- **MINOR**: Incremented for backward-compatible new functionality
- **PATCH**: Incremented for backward-compatible bug fixes

### v0.x Versioning

- During the v0.x phase, minor version increments (e.g., v0.1.0 → v0.2.0) may include breaking changes to the API
- Patch versions (e.g., v0.1.0 → v0.1.1) will only include backward-compatible bug fixes
- The v0.x phase is explicitly considered a development phase where the API may evolve

### v1.0.0 and Beyond

Once v1.0.0 is released, the API will be considered stable, and the project will adhere strictly to semantic versioning:

- Breaking changes will only occur in major version increments (e.g., v1.x.x → v2.0.0)
- New features that don't break backward compatibility will increment the minor version (e.g., v1.0.0 → v1.1.0)
- Bug fixes that don't break backward compatibility will increment the patch version (e.g., v1.0.0 → v1.0.1)

## Pre-release Versions

For significant releases, pre-release versions may be used:

- Alpha: `v0.1.0-alpha.1`, `v0.1.0-alpha.2`, etc.
- Beta: `v0.1.0-beta.1`, `v0.1.0-beta.2`, etc.
- Release Candidate: `v0.1.0-rc.1`, `v0.1.0-rc.2`, etc.

## Criteria for v1.0.0

We will release v1.0.0 when:

1. The API has proven to be functional and well-designed through real-world usage
2. Documentation is comprehensive and clear
3. Test coverage remains high (85%+)
4. No significant changes to the core API are anticipated

## Version Tags

All releases will be tagged in the git repository using the format `v1.2.3` (including the 'v' prefix) to align with Go module conventions.