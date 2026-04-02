# PairAdmin v2.0 Release Checklist

## Pre‑Release
- [x] All unit tests pass (`go test ./internal/...`)
- [x] Integration tests pass (`./scripts/test-integration.sh`)
- [x] End‑to‑end tests pass (`./scripts/test-e2e.sh`)
- [x] No known critical bugs
- [x] Documentation updated (user guide, installation)
- [x] Version number updated in `wails.json`

## Packaging
- [x] macOS binary builds successfully (`.app` bundle)
- [x] Windows binary builds successfully (`.exe`)
- [x] Linux binary builds successfully (executable)
- [ ] Installers tested on clean VMs

## Release
- [x] Create GitHub release with changelog
- [x] Upload all platform binaries
- [ ] Update website/download page
- [ ] Announce on social channels

## Post‑Release
- [ ] Monitor error logs
- [ ] Gather user feedback
- [ ] Plan next iteration

**Note**: Packaging and release steps require a remote Git repository with GitHub Actions enabled. Tag `v2.0.0` has been created locally. Push to remote to trigger automated builds.