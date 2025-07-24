# Weather MCP

A complete Model Context Protocol (MCP) project featuring both a weather server and a client, built in Python. This project demonstrates how to expose weather data and tools to LLMs and interact with them via a custom client.

## Features
- **Weather MCP Server**: Provides weather alerts and forecasts using the National Weather Service API.
- **MCP Client**: Connects to the server, interacts with LLMs (Claude via Anthropic), and demonstrates tool usage.
- **Modern Python project structure** with isolated environments for server and client.

## Directory Structure
```
weather-mcp/
├── weather/
│   ├── weather.py         # MCP server implementation
│   ├── .gitignore         # Server-side ignores
│   └── ...
│   └── mcp-client/
│       ├── client.py      # MCP client implementation
│       ├── .gitignore     # Client-side ignores
│       └── ...
├── README.md
└── ...
```

## Setup

### 1. Clone the Repository
```sh
git clone https://github.com/Rahul-Kaura/weather-mcp.git
cd weather-mcp
```

### 2. Set Up the Server
```sh
cd weather
uv venv
source .venv/bin/activate
uv add "mcp[cli]" httpx
```

### 3. Set Up the Client
```sh
cd mcp-client
uv venv
source .venv/bin/activate
uv add mcp anthropic python-dotenv
```

### 4. Add Your Anthropic API Key
Create a `.env` file in `weather/mcp-client/`:
```
ANTHROPIC_API_KEY=your-anthropic-api-key-here
```
**Never commit your API key!**

## Usage

### Start the Server
```sh
cd weather
source .venv/bin/activate
python weather.py
```

### Start the Client
```sh
cd weather/mcp-client
source .venv/bin/activate
python client.py ../weather.py
```

You’ll see an interactive prompt. Type queries like:
- `What’s the weather in Sacramento?`
- `What are the active weather alerts in Texas?`

Type `quit` to exit.

## Security Notes
- `.env` and `.venv` are git-ignored for safety.
- **Never commit your API keys or secrets.**

## References
- [Model Context Protocol Quickstart](https://modelcontextprotocol.io/quickstart/server)
- [MCP Python SDK on PyPI](https://pypi.org/project/mcp/)

---

Feel free to fork, star, and contribute! 