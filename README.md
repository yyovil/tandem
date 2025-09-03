<p align="center" style="border-radius: 12px; overflow: clip;"><video src="https://github.com/user-attachments/assets/bbfd3fdc-a196-4070-8a7e-c9fad97a322d" width="800" autoplay loop muted></video></p>

<p>
Swarm of AI Agents to assist in a penetration testing engagement given a RoE.md stating Rules of Engagement.
</p>

## Installation

### Quick Install

Install Tandem with a single command using our installation script:

```bash
curl -fsSL https://tandem.codes/install.sh | bash
```

The installation script will:
- Download the latest release for your platform (Linux/macOS, x86_64/arm64)
- Install Tandem to `~/.tandem/bin`
- Add the binary to your shell's PATH
- Verify Docker is installed and running (required for agent tooling)
- Set up default configuration files

### Package Managers

Support for popular package managers is coming soon:
- **Homebrew** (macOS/Linux)
- **Nix** (NixOS/Nix package manager)
- **APT/YUM** (Linux distributions)

*Package manager support will be available by end of day.*

### Verification

After installation, verify Tandem is working:

```bash
tandem --version
```

## Configuration

### Setting up API Keys

Before using Tandem, you need to configure API keys for the AI providers you want to use. 

1. Copy the example environment file:
   ```bash
   cp .example.env .env
   ```

2. Edit the `.env` file and add your API keys for the providers you plan to use:
   ```bash
   GEMINI_API_KEY=your_gemini_api_key_here
   OPENAI_API_KEY=your_openai_api_key_here
   GROQ_API_KEY=your_groq_api_key_here
   OPENROUTER_API_KEY=your_openrouter_api_key_here
   VERTEX_API_KEY=your_vertex_api_key_here
   XAI_API_KEY=your_xai_api_key_here
   ANTHROPIC_API_KEY=your_anthropic_api_key_here
   COPILOT_API_KEY=your_copilot_api_key_here
   ```

   **Note**: You don't need to configure all providers - only add keys for the services you want to use.

### Agent Configuration

Tandem's behavior is controlled by the `.tandem/swarm.json` configuration file, which defines the AI agents and their roles. The default configuration includes several specialized agents:

#### Available Providers
The following AI providers are supported:
- **Groq**: Fast inference for various open-source models
- **Anthropic**: Claude models for reasoning and analysis
- **OpenAI**: GPT models for general-purpose tasks
- **Gemini**: Google's AI models
- **OpenRouter**: Access to multiple AI models through a single API
- **Vertex AI**: Google Cloud AI platform
- **GitHub Copilot**: AI pair programming assistant
- **xAI**: Grok and other xAI models

#### Default Agents

**Orchestrator Agent**
- **Role**: Coordinates and assigns penetration testing tasks to specialized agents
- **Purpose**: Translates user objectives (RoE + chat) into concrete task briefs with suggested tools & techniques; dispatches work to other agents an*d tracks progress
- **Tools**: subagent (delegates tasks to other agents)

**Reconnoiter Agent**
- **Role**: Seasoned OffSec PEN-300 certified penetration tester with extensive experience in reconnaissance
- **Purpose**: Performs reconnaissance (network/service enumeration, OSINT, surface mapping) to build target knowledge for later phases
- **Tools**: terminal (Kali Linux CLI tooling)

**Vulnerability Scanner Agent**
- **Role**: Vulnerability assessment specialist
- **Purpose**: Runs targeted scans to identify, categorize, and prioritize vulnerabilities discovered during reconnaissance
- **Tools**: terminal (Kali Linux CLI tooling)

**Exploiter Agent**
- **Role**: Exploitation specialist
- **Purpose**: Researches viable exploits for identified vulnerabilities and executes them to gain footholds / escalate access within the allowed RoE boundaries
- **Tools**: terminal (Kali Linux CLI tooling)

**Reporter Agent**
- **Role**: Reporting & analysis specialist
- **Purpose**: Synthesizes findings from all phases into objective, business-impact focused reporting with actionable remediation recommendations

## Usage

After configuring your API keys and agent settings:

1. **Set up your engagement context**: Create a `RoE.md` file in your working directory containing the Rules of Engagement for your penetration testing engagement.

2. **Run Tandem**: Start the TUI interface to interact with your AI agent swarm:
   ```shell
   tandem
   ```

3. **Interact with agents**: Use the interface to communicate with specialized agents for different phases of your penetration testing workflow.