---
name: repository-guidance
description: 'Create or update concise repository guidance for AI coding agents. Use when bootstrapping AGENTS.md or copilot-instructions.md, documenting project architecture and commands, or improving agent productivity without duplicating existing docs.'
argument-hint: 'What project area or workflow should the agent guidance cover?'
user-invocable: true
disable-model-invocation: true
---

# Repository Guidance

Create or improve workspace-scoped instructions that help coding agents become productive quickly. The output is an appropriate agent-instructions file, usually `AGENTS.md` at the repository root, plus focused validation.

## Procedure

1. Identify the repository root and search for existing customization files: `AGENTS.md`, `.github/copilot-instructions.md`, `CLAUDE.md`, and other agent-rule files. Preserve valuable existing guidance; update rather than duplicate it.
2. Read the nearest project documentation and build metadata. Verify commands from files such as `README.md`, `Makefile`, `package.json`, `go.mod`, or equivalent configuration instead of guessing.
3. Inspect the smallest useful set of source files, tests, and generated-asset declarations to establish:
   - package or component ownership and data flow;
   - build, test, generation, and run commands;
   - project-specific conventions and development prerequisites;
   - runtime data, generated files, concurrency, integration, or environment pitfalls.
4. Choose the smallest customization primitive. Prefer one root `AGENTS.md` for repository-wide guidance; add scoped instruction files only when a distinct area needs different rules.
5. Draft concise, imperative guidance. Include links to existing documentation rather than copying it. Record only facts that are useful to an agent and difficult to infer safely from the code.
6. Validate the customization file: confirm its location and frontmatter if present, check that referenced files exist, and run the narrowest relevant project test or check. Do not change application code as part of this workflow.
7. Review the draft for ambiguity and duplication. Ask the user about any unresolved scope or convention question, then summarize the resulting file and suggest one or two related customizations.

## Quality Criteria

- Every command and prerequisite is verified against the repository.
- Instructions name the owning files or packages for important behavior.
- Existing documentation is linked, not reproduced.
- Generated files, persistent data, and environment-specific artifacts are clearly distinguished from source.
- Warnings describe actionable failure modes, not speculative design commentary.
- The change is limited to chat customization files and leaves unrelated worktree changes intact.

## Completion Report

Report the added or modified customization files in a small table, explain why each helps agents, list the validation performed, and mention any remaining uncertainty. For iterative improvement, use `/chronicle improve` to identify recurring friction from later sessions.
