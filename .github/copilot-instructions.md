# Coding Preferences

Follow these personal preferences while writing code when asked:

- use ```any``` instead of ```interface{}```
- lookup the go.mod file in the workspace to understand the pkgs used in the project and stick your code suggestions around them. do not bring any other dependency by yourself unless explicitly approved by the user.
- use `errors.Wrap(err, "message")` for error handling

# Project description

Tandem is a terminal app for interfacing with a 
Swarm of AI Agents that work in tandem to assist during a Penetration Testing engagement. It should be given a Rules of Engagement(RoE) file to build the context boundary and then upon given the tasks it will execute them and accomplish the intended end goal as mentioned in the RoE by using the help of its Agents.

