# PairAdmin Project Checklist Usage Guide

## Overview
The `PROJECT_CHECKLIST.json` file is designed to help an LLM agent track progress across multiple sessions while implementing PairAdmin v2.0. It contains all 50 atomic tasks and 12 QA checkpoints with their dependencies, statuses, and verification criteria.

## File Structure
The JSON file has the following top‑level keys:

- `project`: Metadata about the project (name, version, current phase, status).
- `phases`: Four development phases (Foundation, Terminal Integration, Windows & Polish, Hardening & Launch).
- `tasks`: Array of all 50 tasks with detailed information.
- `qa_checkpoints`: Array of 12 QA checkpoints that correspond to milestones.
- `notes`: Important usage notes.

## Task Object Format
Each task has:
```json
{
  "id": 1,
  "phase": "1",
  "title": "Initialize Go module",
  "description": "Create go.mod with Go 1.21+",
  "status": "pending",
  "dependencies": [],
  "artifacts": ["go.mod"],
  "verification": ["go mod tidy succeeds", "go version shows 1.21+"]
}
```

**Status values**:
- `pending`: Not yet started.
- `in_progress`: Currently being worked on.
- `completed`: Successfully finished.
- `blocked`: Cannot proceed due to external dependency or bug.

## Checkpoint Object Format
Each checkpoint has:
```json
{
  "id": "CP1",
  "name": "Milestone 1: Project Setup Complete",
  "tasks_covered": [1, 2, 3, 4, 5],
  "status": "pending",
  "validation_criteria": ["go.mod exists", "wails.json exists", ...],
  "required_artifacts": ["go.mod", "wails.json", ...]
}
```

**Checkpoint status values**:
- `pending`: Not yet validated.
- `in_progress`: Validation in progress.
- `passed`: All validation criteria met.
- `failed`: One or more criteria failed.

## How to Use the Checklist

### 0. Helper Scripts (Optional)

Two helper scripts are provided to make checklist management easier:

1. **`update_checklist.py`**: Python script with command‑line interface.
2. **`checklist.sh`**: Bash wrapper that calls the Python script.

**Common commands**:
```bash
# Show progress summary
./checklist.sh summary

# Update task status
./checklist.sh task 1 completed

# Update checkpoint status  
./checklist.sh checkpoint CP1 passed

# List pending tasks
./checklist.sh list tasks --status pending

# List all checkpoints
./checklist.sh list checkpoints
```

You can use these scripts or edit the JSON file directly. The scripts automatically update the `last_updated` timestamp.

### 1. At the Start of a Session
When you begin a new session or continue work:

1. **Read the checklist**:
   ```bash
   cat PROJECT_CHECKLIST.json
   ```
   or use your file reading tool.

2. **Determine current state**:
   - Check `project.current_phase` to see which phase is active.
   - Look for tasks with `status: "in_progress"` (there should be at most one).
   - If no task is `in_progress`, find the first `pending` task whose dependencies are all `completed`.

3. **Plan your work**:
   - Select a task to work on.
   - Update its status to `in_progress`.
   - Update `project.last_updated` timestamp.

### 2. While Working on a Task
1. **Reference the task details**:
   - Look at `description`, `artifacts`, `verification`.
   - Check `dependencies` to ensure they are `completed`.

2. **Follow the implementation plan**:
   - Refer to `IMPLEMENTATION_PLAN.md` for detailed steps.
   - Refer to `TASK_EXAMPLE.md` for task specification format.

### 3. After Completing a Task
1. **Verify the task**:
   - Run each verification step listed in the task.
   - Ensure all artifacts exist and are correct.

2. **Update the checklist**:
   - Set task `status` to `completed`.
   - Add a `completed_at` timestamp (optional, you can add field).
   - Update `project.last_updated`.

3. **Check for checkpoint completion**:
   - If all tasks in a checkpoint are `completed`, you may validate the checkpoint.
   - Follow validation criteria in `QA_CHECKPOINTS.md`.
   - Update checkpoint `status` to `passed` if successful.

### 4. When Blocked
If a task cannot proceed due to external issues:
- Set `status` to `blocked`.
- Add a `blocked_reason` field (optional).
- Move to another task whose dependencies are satisfied.

### 5. Between Sessions
Before ending a session:
- Ensure at most one task is `in_progress`.
- Update `project.last_updated`.
- Consider leaving a note about what to do next in `project.notes` (optional).

## Example Workflow

**Session 1**:
1. Read checklist: all tasks pending.
2. Start task 1: update `status` to `in_progress`.
3. Implement task 1 (create `go.mod`).
4. Verify: `go mod tidy` succeeds.
5. Update task 1 `status` to `completed`.
6. Start task 2: update `status` to `in_progress`.
7. End session.

**Session 2**:
1. Read checklist: task 2 `in_progress`, task 1 `completed`.
2. Continue task 2 (scaffold Wails project).
3. Verify: `wails build` succeeds.
4. Update task 2 `status` to `completed`.
5. Check checkpoint CP1: tasks 1‑5 all completed? No, tasks 3‑5 pending.
6. Start task 3.
7. ...

## JSON Manipulation Examples

### Reading the file (in LLM context)
You can parse the JSON mentally or use tools if available. The structure is straightforward.

### Updating a task status
Assume you have the JSON loaded as a JavaScript object (or you can edit the text directly):

```javascript
// Pseudo‑code for updating task 1 to completed
const checklist = JSON.parse(fileContent);
const task = checklist.tasks.find(t => t.id === 1);
task.status = "completed";
task.completed_at = new Date().toISOString(); // optional
checklist.project.last_updated = new Date().toISOString();
const updatedContent = JSON.stringify(checklist, null, 2);
// Write back to file
```

If you cannot run JavaScript, edit the JSON text directly:
- Find the task object.
- Change `"status": "pending"` to `"status": "completed"`.
- Update the `last_updated` field.

### Adding custom fields
You may add fields to tasks or checkpoints as needed (e.g., `completed_at`, `blocked_reason`, `notes`). Maintain the overall structure.

## Best Practices

1. **One task at a time**: Only have one task with `status: "in_progress"`.
2. **Validate dependencies**: Never start a task unless all its dependencies are `completed`.
3. **Update timestamps**: Always update `project.last_updated` when modifying the file.
4. **Check checkpoints**: After completing a group of tasks, validate the corresponding checkpoint using `QA_CHECKPOINTS.md`.
5. **Keep artifacts**: Save the artifacts listed for each task; they may be needed later.
6. **Document blockers**: If a task is blocked, note why so the next session can address it.

## Integration with Other Documents

- **`IMPLEMENTATION_PLAN.md`**: Contains detailed descriptions of each task.
- **`QA_CHECKPOINTS.md`**: Provides validation criteria for each milestone.
- **`COMPREHENSIVE_QA_PLAN.md`**: Full QA strategy for final validation.
- **`TASK_EXAMPLE.md`**: Example of how to specify a single task.
- **`update_checklist.py`** / **`checklist.sh`**: Helper scripts for managing the checklist.

## Recovery from Corruption
If the JSON file becomes malformed:
- Refer to the versions in this document to reconstruct.
- Use the task list from `IMPLEMENTATION_PLAN.md` as backup.
- The LLM can regenerate the checklist from the implementation plan if needed.

## Notes for Frontier Models
- You have a large context window; you can load the entire checklist and related documents.
- Use the checklist to maintain state across sessions; it's your "memory".
- Be thorough in verification; don't mark tasks `completed` until all criteria are met.
- When in doubt, consult the PRD (`/home/sblanken/download/PairAdmin_PRD_v2.0.md`).

---

*This checklist is your primary tool for tracking progress. Update it faithfully after each significant step.*