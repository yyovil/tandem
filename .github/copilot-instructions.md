# Coding Preferences

Follow these personal preferences while writing code when asked:

- use ```any``` instead of ```interface{}```
- lookup the go.mod file in the workspace to understand the pkgs used in the project and stick your code suggestions around them. do not bring any new dependency by yourself unless explicitly approved | asked by the user.
- when asked to generate stubs for implementing interfaces, only provide the method signatures with empty bodies. 
  <md>
  <br/>

  **Example**

  **user:** generate a stub implementing #BaseTool interface for #DockerCli.

  **assistant**: here's the stub implementing the #BaseTool interface for the DockerCli:
    ```go
    func (d *DockerCli) Info() tools.ToolInfo {}
    ```

  </md>


# Project description

[Tandem](https://github.com/yaydraco/tandem) is a terminal app for interfacing with a 
Swarm of AI Agents that work in tandem to assist during a Penetration Testing engagement. It should be given a Rules of Engagement(RoE) file to build the context boundary and then upon given the tasks it will execute them and accomplish the intended end goal as mentioned in the RoE by using the help of its Agents.
List of Penetration testing related AI Agents at disposal:
- **Orchestrator**: Point of Contact for the users to interact with the Swarm of AI Agents. user will have to explain the task to the orchestrator to complete. 
- **Reconnoier**: Agent that performs reconnaissance tasks.
- **Vulnerability Scanner**: Agent that scans for vulnerabilities.
- **Exploiter**: Agent that exploits the vulnerabilities found by the Vulnerability Scanner.
- **Reporter**: Agent that generates reports based on the findings of the other agents.

# Response Preferences in Ask mode.

- response for "CONCEPTUAL" queries should be answered colloquially, short enough just to get the gist of it. Please don't use vague and off-topic examples to prove your point. Keep the explanation contextual, relevant and fully contained.
