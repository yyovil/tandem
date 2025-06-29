## Defaults
- database is created at ~/.tandem/data.
- $PWD is the bind mount for the swarm's container.
- configuration at ~/.config/.tandem/
- .tandem/RoE.md contains the context about the engagement. its expected to be present in every dir where the engagement data is stored.
> **NOTE**: *These settings are applied as defaults, unless specified.*

## Decisions
- there won't be themes. there would be a theme to avoid lots of code refactoring that adheres to the opencode's api but other than that, we don't need themes here.
- don't worry about config for now.
- instead of CoderAgent we would have something like ochestrator as Point of Contact for the users. Tandem tui is basically the interface for the user to interact with the Swarm of AI Agents.
- instead of AgentTool that allows the CoderAgent to spawn subagents for tasks parallelization and save the context window, we would have something like Delegate | Assign for the orchestrator.
- we are going to need the edit tool. project manager is going to need it.
> **NOTE**: *These are some design decisions that tandem follows*.


## 2nd Iteration
- replace this pkg https://github.com/ncruces/go-sqlite3 with https://github.com/mattn/go-sqlite3.
- we can use go:generate for querying https://models.dev for models. then we can generate json schema and the internal/models/*.go.
- opencode uses model specific base prompts. now should that be the case with the tandem could only be decided once any differences are observed during the engagement time.
> **NOTE**: *Try these suggestions out in the 2nd iteration. Some of them might feel stupid to you in the retrospect.*