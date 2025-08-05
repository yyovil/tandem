## Defaults
- database could be found at ~/.tandem/data.
- $PWD is the bind mount for the swarm's container.
- configuration could be found at ~/.config/.tandem/
- .tandem/RoE.md contains the context about the engagement. its expected to be present in every dir where the engagement data is stored.
> **NOTE**: *These settings are applied as defaults, unless specified.*

## Decisions
> **NOTE**: *These are some design decisions that tandem follows*.


## Experiments
- we can use **go:generate** for querying https://models.dev for models. then we can generate json schema and the internal/models/*.go.
- opencode uses model specific base prompts. now should that be the case with the tandem could only be decided once any differences are observed during the engagement time.
- what if we ***allow the orchestrator to spawn subagents*** during the runtime?
> **NOTE**: *These experiments are kind of like abalation, QoL improvements etc.*

## Wish
- agent's have a tool that they can use to notify the user when working in the background.

## Rough
- responseSchema and mimeType are provider specific because not every provider supports it. then they are passed while creating the provider client and we gotta do the prop drilling from NewAgent -> createAgentProvider.


