# Question Answering MCP for Cursor

This is a simple Model Context Protocol (MCP) server that provides a question-answering tool for Cursor.

## Features

- Simple HTTP server that follows the MCP protocol
- Exposes a tool to ask questions and get answers
- Works with Cursor's MCP integration

## Setup

1. Make sure you have Go installed on your system
2. Clone this repository
3. Run the server:
   ```
   go run main.go
   ```

## Usage with Cursor

### Option 1: Using the included configuration

This project includes a `.cursor/mcp.json` file that Cursor will automatically detect when you open this folder as a workspace. The MCP server will be available to Cursor automatically.

### Option 2: Manual configuration in Cursor

1. Open Cursor
2. Go to Settings > Features > MCP
3. Click "+ Add New MCP Server"
4. Configure as follows:
   - Name: Question Answering MCP
   - Type: stdio
   - Command: `go run /path/to/this/repo/main.go`
5. Click "Add Server"

## Using the MCP in Cursor

Once configured, the MCP server will provide the `ask_question` tool to Cursor's AI assistant. You can use it by:

1. Opening Cursor's AI assistant (Composer)
2. Asking it to use the question-answering tool
3. The AI will call the tool and display the response

Example prompt: "Use the question-answering tool to ask: What is the capital of France?"

## Extending the MCP

To add more tools to this MCP server:

1. Add new tool definitions in the `handleMCPTools` function
2. Implement the tool's logic in the `handleMCPExecute` function
3. Restart the server

## License

MIT 