# PairAdmin v2.0 Comprehensive QA Plan

## 1. Overview and Objectives
This document defines a comprehensive quality assurance plan for PairAdmin v2.0. The plan is designed to be executed by a frontier model with large context window capabilities, ensuring thorough validation of all functional, non‑functional, security, and cross‑platform requirements before release.

**Objectives**:
- Validate all user stories and acceptance criteria from the PRD.
- Ensure cross‑platform compatibility (Linux, macOS, Windows).
- Verify security and privacy safeguards.
- Assess performance, usability, and reliability.
- Provide actionable feedback for final improvements.

**Scope**: Covers the entire PairAdmin v2.0 application as defined in the PRD, including terminal integrations, LLM gateway, UI, settings, and security features.

**Out of scope**: Modifications to upstream dependencies (tmux, terminal emulators, LLM APIs).

## 2. Testing Environments and Setup

### 2.1 Hardware/Software Matrix
| Platform | OS Version | Terminal Emulators | LLM Providers | Notes |
|----------|------------|-------------------|---------------|-------|
| **Linux** | Ubuntu 22.04 LTS | tmux 3.2a, GNOME Terminal 3.44, Konsole 22.12 | OpenAI (GPT‑4), Ollama (llama2) | AT‑SPI2 enabled |
| **Linux** | Fedora 38 | tmux 3.3, GNOME Terminal 3.48 | Anthropic (Claude), Local llama.cpp | |
| **macOS** | macOS Ventura 13.5 | tmux 3.3, Terminal.app | OpenAI, Anthropic, Ollama | Accessibility permission granted |
| **macOS** | macOS Sonoma 14.2 | Terminal.app, iTerm2 (optional) | OpenAI, Local LM Studio | |
| **Windows** | Windows 11 22H2 | PuTTY 0.78, Windows Terminal (optional) | OpenAI, Ollama (via WSL) | UI Automation enabled |
| **Windows** | Windows 10 21H2 | PuTTY 0.77 | OpenAI | |

### 2.2 Test Data
- **Terminal content samples**:
  - Common commands (`ls`, `cd`, `cat`, `grep`)
  - Error messages (permission denied, file not found, network errors)
  - Multi‑line output (e.g., `docker ps`, `kubectl get pods`)
  - Sensitive data patterns (AWS keys, passwords, private keys)
- **LLM test prompts**:
  - "What does this error mean?" (with error output)
  - "Suggest a command to restart nginx"
  - "How do I check disk usage?"
  - "Explain this output: <multi‑line output>"

### 2.3 Pre‑requisites
- PairAdmin built from source (or installed via package) on each platform.
- LLM API keys for cloud providers (OpenAI, Anthropic) stored in OS keychain.
- Ollama installed and running with at least one model (e.g., `llama2`).
- Terminal emulators installed and configured.
- Test user accounts with appropriate permissions.

## 3. Test Categories

### 3.1 Unit Testing
**Goal**: Verify individual components in isolation.

**Components to test**:
- `internal/session/` – Session creation, retrieval, closure.
- `internal/context/` – Circular buffer, token counting, truncation.
- `internal/llm/providers/` – Each provider (OpenAI, Anthropic, Ollama) with mocked HTTP.
- `internal/security/filter.go` – Regex pattern matching and redaction.
- `internal/terminal/tmux/adapter.go` – Parsing of `tmux` output (mocked `exec.Command`).

**Tools**: Go’s `testing` package, `testify` for assertions, `gomock` for interfaces.

**Success Criteria**: >70% line coverage for critical paths; all unit tests pass.

### 3.2 Integration Testing
**Goal**: Verify interactions between components.

**Integration points**:
1. **Terminal detection → Session manager**: New terminal detected creates a session.
2. **Context manager → LLM gateway**: Terminal content filtered and sent to LLM.
3. **Clipboard manager → System clipboard**: Copy/paste works across platforms.
4. **Settings manager → LLM gateway**: Provider change takes effect immediately.
5. **Security filter → Context manager**: Redacted content never leaves machine (for cloud LLMs).

**Tools**: Integration tests with limited mocking (e.g., real tmux session, mocked LLM).

**Success Criteria**: All integration tests pass; components cooperate without data corruption.

### 3.3 System Testing (End‑to‑End)
**Goal**: Validate complete user workflows.

**Test scenarios** (based on user stories):
1. **US‑001**: Launch PairAdmin, open a terminal, type commands; verify PairAdmin captures output within 500ms.
2. **US‑002**: Ask AI about terminal content; verify response references the content and suggests commands.
3. **US‑003**: Generate multiple AI suggestions; verify sidebar shows them in reverse chronological order.
4. **US‑004**: Open two terminals; verify separate tabs, isolated chat histories.
5. **US‑005**: Click “Copy to Terminal”; verify command appears in clipboard and can be pasted.
6. **US‑006**: Open settings, change LLM provider, add custom prompt; verify AI behavior changes.

**Tools**: Semi‑automated scripts that simulate user input (using UI automation libraries like `selenium` for webviews, `robotgo` for native UI). Manual validation for complex interactions.

**Success Criteria**: All acceptance criteria pass on at least one platform per terminal type.

### 3.4 Cross‑Platform Testing
**Goal**: Ensure consistent behavior across Linux, macOS, Windows.

**Platform‑specific tests**:
- **Linux**: AT‑SPI2 accessibility, GNOME Terminal/Konsole capture, `xclip`/`xsel` clipboard.
- **macOS**: Accessibility API permission flow, Terminal.app capture, macOS keychain.
- **Windows**: UI Automation for PuTTY, Windows Credential Manager, clipboard API.

**Success Criteria**: All core features work on each platform; platform‑specific adapters handle errors gracefully.

### 3.5 Security Testing
**Goal**: Verify sensitive data protection and secure operation.

**Test areas**:
1. **Data redaction**: Terminal content containing AWS keys, passwords, private keys is redacted before being sent to cloud LLM.
2. **Local LLM**: When Ollama/LM Studio is selected, no network traffic leaves the machine (verified with Wireshark/tcpdump).
3. **Credential storage**: API keys stored in OS keychain, not plain config files.
4. **Audit logging**: Logs contain redacted events; no sensitive data written to disk.
5. **Clipboard hygiene**: PairAdmin‑copied commands cleared from clipboard after timeout (configurable).
6. **Permission handling**: macOS Accessibility permission request works; fallback guidance provided.

**Tools**: Static analysis (`gosec`), dynamic analysis (network monitoring, file monitoring), manual review.

**Success Criteria**: No sensitive data leaked; all security constraints from PRD section 6 satisfied.

### 3.6 Performance Testing
**Goal**: Ensure responsiveness and resource efficiency.

**Metrics**:
- **Terminal capture latency**: <500ms from terminal change to UI update.
- **LLM response time**: <30 seconds for typical queries (depends on provider).
- **CPU usage**: <5% when idle (polling terminals).
- **Memory usage**: <200 MB for typical session; no leaks over 24‑hour stress test.
- **Startup time**: <3 seconds from launch to UI ready.

**Tools**: Go’s `pprof` for profiling, custom benchmarks, `stress` scripts.

**Success Criteria**: Meets or exceeds performance targets; no memory leaks.

### 3.7 Usability Testing
**Goal**: Evaluate user interface against wireframe specifications and intuitive use.

**Checklist**:
- Three‑column layout matches wireframe (widths, spacing, colors).
- Dark/light theme switching works.
- Tab switching is fast and clear.
- “Copy to Terminal” button is prominent and intuitive.
- Settings dialog is well‑organized and understandable.
- Error messages are helpful and guide recovery.

**Method**: Heuristic evaluation by QA engineer; small user study with target personas (if possible).

**Success Criteria**: UI conforms to wireframes; no major usability issues.

### 3.8 Compatibility Testing
**Goal**: Verify compatibility with supported terminal emulators and LLM providers.

**Terminal matrix** (as per PRD):
- tmux (Linux, macOS, Windows via WSL)
- GNOME Terminal (Linux)
- Konsole (Linux)
- Terminal.app (macOS)
- PuTTY (Windows)

**LLM provider matrix**:
- OpenAI (GPT‑4, GPT‑3.5‑turbo)
- Anthropic (Claude‑3‑Haiku, Claude‑3‑Sonnet)
- Ollama (llama2, mistral, etc.)
- LM Studio (local)
- llama.cpp (self‑hosted)

**Success Criteria**: Each combination works for basic capture and chat; edge cases documented.

## 4. Detailed Test Cases

### 4.1 US‑001: Automatic Terminal Context Capture
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑001‑01 | Capture tmux session | 1. Start tmux session.<br>2. Launch PairAdmin.<br>3. Type `ls` in tmux.<br>4. Wait 500ms. | PairAdmin tab shows `ls` output in terminal preview. |
| TC‑001‑02 | Capture GNOME Terminal | 1. Open GNOME Terminal.<br>2. PairAdmin detects it.<br>3. Type `pwd`.<br>4. Wait 500ms. | Output appears in preview. |
| TC‑001‑03 | Multiple terminals isolation | 1. Open two tmux sessions.<br>2. PairAdmin shows two tabs.<br>3. Type different commands in each.<br>4. Switch tabs. | Each tab shows only its own terminal content. |
| TC‑001‑04 | Context limit prioritization | 1. Generate >1000 lines of output (e.g., `yes`).<br>2. Check PairAdmin preview. | Preview shows most recent lines (within context limit). |

### 4.2 US‑002: Contextual AI Assistance
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑002‑01 | AI references terminal context | 1. Terminal shows "Permission denied".<br>2. Ask "What went wrong?"<br>3. Submit. | AI response mentions "Permission denied" and suggests `chmod` or `sudo`. |
| TC‑002‑02 | Command suggestion with button | 1. Ask "How do I list files?"<br>2. AI responds. | Response includes code block with `ls -la` and a "Copy to Terminal" button. |
| TC‑002‑03 | Copy to Terminal button copies | 1. Click "Copy to Terminal" button.<br>2. Paste into terminal. | Command appears in terminal. |
| TC‑002‑04 | Auto‑paste (if configured) | 1. Enable auto‑paste in settings.<br>2. Click button. | Command is pasted directly into focused terminal. |

### 4.3 US‑003: Command History Sidebar
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑003‑01 | Sidebar shows commands | 1. Generate 3 AI‑suggested commands.<br>2. Look at sidebar. | Three command cards appear, newest first. |
| TC‑003‑02 | Click to copy | 1. Click a command card.<br>2. Paste elsewhere. | Command text is in clipboard. |
| TC‑003‑03 | Hover shows context | 1. Hover over command card.<br>2. Wait for tooltip. | Tooltip displays original question that generated the command. |
| TC‑003‑04 | Smooth scrolling | 1. Generate 20 commands.<br>2. Scroll sidebar. | Scrolling is smooth; no jumps or freezes. |

### 4.4 US‑004: Multi‑Terminal Tab Management
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑004‑01 | Auto‑tab creation | 1. Open a new terminal (tmux/GNOME/PuTTY).<br>2. Wait for detection. | New tab appears in PairAdmin. |
| TC‑004‑02 | Isolated chat per tab | 1. Tab A: ask "What is pwd?"<br>2. Switch to Tab B.<br>3. Chat history is empty. | Chat histories are separate. |
| TC‑004‑03 | Terminal closure prompt | 1. Close a terminal window.<br>2. PairAdmin detects closure. | User is prompted to close or archive the tab. |

### 4.5 US‑005: One‑Click Command Execution
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑005‑01 | Clipboard injection | 1. Click "Copy to Terminal".<br>2. Check clipboard content. | Clipboard contains the command. |
| TC‑005‑02 | Hotkey copies last command | 1. Configure hotkey Ctrl+Shift+C.<br>2. Press hotkey. | Most recent command copied to clipboard. |
| TC‑005‑03 | Hotkey focuses terminal | 1. Press hotkey (if configured to focus).<br>2. Check active window. | Target terminal becomes focused. |
| TC‑005‑04 | Capture after execution | 1. Execute a command (e.g., `echo "test"`).<br>2. Wait for output. | PairAdmin captures the new output. |

### 4.6 US‑006: LLM Configuration
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑006‑01 | Provider selection | 1. Open settings.<br>2. Select "Anthropic".<br>3. Enter API key.<br>4. Save. | Provider changes; test connection succeeds. |
| TC‑006‑02 | Custom prompt extension | 1. Add "Always respond in Spanish." to custom prompt.<br>2. Ask a question. | AI responds in Spanish. |
| TC‑006‑03 | Local LLM (Ollama) | 1. Select Ollama provider.<br>2. Ask question. | Response from local model; no external network traffic. |

### 4.7 Security Test Cases
| Test Case ID | Description | Steps | Expected Result |
|--------------|-------------|-------|-----------------|
| TC‑SEC‑01 | Password redaction | 1. Terminal shows `password: secret123`.<br>2. Ask AI about terminal.<br>3. Inspect network request. | Redacted: `password: [REDACTED]`. |
| TC‑SEC‑02 | AWS key redaction | 1. Terminal contains `AKIA...`.<br>2. Send to cloud LLM. | Key replaced with `[AWS_KEY_REDACTED]`. |
| TC‑SEC‑03 | Local LLM no exfiltration | 1. Select Ollama.<br>2. Monitor network with `tcpdump`. | No traffic to external IPs. |
| TC‑SEC‑04 | Audit log redaction | 1. Perform several actions.<br>2. Check audit log file. | Log entries contain redacted content. |
| TC‑SEC‑05 | Keychain storage | 1. Save API key.<br>2. Inspect config file. | Config file does not contain plain‑text key. |

### 4.8 Cross‑Platform Test Cases
| Test Case ID | Platform | Description | Expected Result |
|--------------|----------|-------------|-----------------|
| TC‑CP‑01 | macOS | Accessibility permission flow | If not granted, dialog guides user to System Preferences. |
| TC‑CP‑02 | Linux | AT‑SPI2 detection | GNOME Terminal detected when accessibility enabled. |
| TC‑CP‑03 | Windows | PuTTY UI Automation | PuTTY window text captured (or fallback OCR). |
| TC‑CP‑04 | All | Clipboard copy/paste | Command copied to clipboard works in native terminal. |

## 5. Automation Strategy

### 5.1 Automated Tests
- **Unit and integration tests**: Run in CI on every commit (`go test ./...`).
- **Component tests**: Use `testify` and mocks to test adapters in isolation.
- **E2E smoke tests**: Scripts that launch PairAdmin, simulate basic user interactions (using `robotgo` or similar UI automation). These may be platform‑specific and run nightly.

### 5.2 Manual Tests
- **Usability and visual design**: Manual inspection against wireframes.
- **Complex workflows**: Multi‑terminal scenarios, permission flows.
- **Platform‑specific edge cases**: Requires physical/virtual machines with different OS versions.

### 5.3 Test Data Management
- **Sensitive data**: Use synthetic patterns (e.g., `AKIA0000000000000000`) for testing redaction.
- **Terminal output**: Scripts to generate realistic command outputs (`ls`, `docker ps`, error messages).
- **LLM responses**: Mock LLM provider for reproducible tests.

## 6. Execution Schedule

### Phase A: Pre‑integration (Weeks 1‑12)
- Run unit and integration tests after each atomic task (as per QA checkpoints).
- Validate each milestone using the checkpoint criteria.

### Phase B: System testing (Weeks 13‑18)
- Execute system test cases for completed features.
- Cross‑platform testing on Linux, macOS, Windows VMs.
- Security testing of data redaction and local LLM.

### Phase C: Performance and usability (Weeks 19‑22)
- Performance profiling and optimization.
- Usability review against wireframes.
- Compatibility testing with terminal/LLM matrix.

### Phase D: Final validation (Weeks 23‑24)
- Full test suite execution on all target platforms.
- Security audit by dedicated reviewer.
- User acceptance testing with sample personas.

## 7. Reporting and Sign‑off

### 7.1 Defect Tracking
- Use issue tracker (GitHub Issues) with labels: `bug`, `security`, `platform‑specific`.
- Severity levels: Critical, High, Medium, Low.
- Each defect must include steps to reproduce, expected/actual results, environment details.

### 7.2 Test Reports
- Daily test execution summary (automated tests pass/fail).
- Weekly test status report covering manual testing progress.
- Final test report before release, including:
  - Test coverage metrics.
  - Defect summary (open/closed).
  - Performance benchmarks.
  - Security assessment.

### 7.3 Release Criteria
**Must‑have** (blocking release):
- All critical and high‑severity defects resolved.
- All acceptance criteria (US‑001 through US‑006) pass on at least one supported platform per terminal type.
- Security tests pass (no data exfiltration, proper redaction).
- Application builds and runs on all three target platforms.

**Nice‑to‑have** (can defer to patch release):
- Performance targets met on all platforms.
- All terminal emulators in matrix working perfectly.
- Usability issues minor and documented.

### 7.4 Sign‑off Process
1. **Development lead** signs off that all features are implemented per PRD.
2. **QA lead** signs off that all test cases pass and release criteria are met.
3. **Security reviewer** signs off that security constraints are satisfied.
4. **Product owner** signs off that the product meets user needs.

## 8. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Terminal capture unreliable on certain platforms | High | Implement fallback mechanisms (e.g., OCR for PuTTY); document limitations. |
| LLM API changes break providers | Medium | Use version‑pinned API clients; monitor provider status pages. |
| Performance degradation with many terminals | Medium | Optimize polling; add configurable polling interval. |
| Accessibility permissions denied by user | Medium | Clear guidance; degrade gracefully (manual paste mode). |
| Cross‑platform packaging issues | High | Use Wails packaging system; test on clean VMs early. |

## 9. Appendices

### Appendix A: Test Environment Setup Scripts
- `scripts/setup‑test‑linux.sh`: Installs tmux, GNOME Terminal, Konsole, AT‑SPI2, Ollama.
- `scripts/setup‑test‑macos.sh`: Guides through Accessibility permission, installs tmux, Ollama.
- `scripts/setup‑test‑windows.ps1`: Installs PuTTY, configures UI Automation, installs Ollama (via WSL).

### Appendix B: Acceptance Criteria Traceability Matrix
Mapping of each acceptance criterion to test cases (to be maintained as test cases evolve).

### Appendix C: Performance Benchmarks
Detailed performance targets and measurement methodology.

---

*This QA plan is a living document and should be updated as the product evolves. Final validation will not cut corners; every aspect will be thoroughly tested before release.*