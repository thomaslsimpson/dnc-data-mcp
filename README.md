# Model Context Protocol (MCP) Service for DNC 

This MCP provides a view into DNC data that will help Agents use the information.

MCP information and documentation can be found here: https://github.com/modelcontextprotocol

Documentation on writing a MCP server can be found here: https://modelcontextprotocol.io/quickstart/server


## Implementation Plan

The steps we will take:
1. Get an MCP service up and running locally and test it with Ollama
  - this should be a simple service that returns a random number
  - the idea is to make a simple service and test the process of using it with Ollama
  - we will also need to write a simple Agent to talk to Ollama running locally and use that to test the service
2. Next we will set up a read only SQL connection to our database using our local credentials
3. Then we will offer a generic MCP feature to use the data
  - this will allow questions about the database itself and will form queries
  - we will add a client test here to ensure that this section will work with Ollama
4. We will build a data dictionary that we can use for reference based on the existing databse tables
5. Then we will add an additional layer to the service which provides more service information for Agents


We will use the local Ollama installation with the model "llama3.2:latest" for testing.


## Information for reference later

Questions we would like to be able to ask the service
- who are our partners?
- which partner made the most money last month?
- which source tags does partner X use?
- which source tags went up in TQ last month?
- which partners have been on the YER the most in the last 6 months?
- what traffic sources does partner X use?


