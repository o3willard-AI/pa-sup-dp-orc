# PairAdmin v2.0 Release Checklist

## Pre‑Release
- [ ] All unit tests pass (`go test ./internal/...`)
- [ ] Integration tests pass (`./scripts/test-integration.sh`)
- [ ] End‑to‑end tests pass (`./scripts/test-e2e.sh`)
- [ ] No known critical bugs
- [ ] Documentation updated (user guide, installation)
- [ ] Version number updated in `wails.json`

## Packaging
- [ ] macOS `.dmg` builds successfully
- [ ] Windows `.msi` builds successfully
- [ ] Linux packages (`.deb`, `.rpm`, `.AppImage`) build successfully
- [ ] Installers tested on clean VMs

## Release
- [ ] Create GitHub release with changelog
- [ ] Upload all platform binaries
- [ ] Update website/download page
- [ ] Announce on social channels

## Post‑Release
- [ ] Monitor error logs
- [ ] Gather user feedback
- [ ] Plan next iteration