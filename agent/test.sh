#!/bin/bash

# Build the agent
go build -o agent

# Test questions
echo "Testing 'who are our partners?'"
./agent "who are our partners?"

echo -e "\nTesting 'which partner made the most money last month?'"
./agent "which partner made the most money last month?"

echo -e "\nTesting 'which source tags does DNC use?'"
./agent "which source tags does DNC use?"

echo -e "\nTesting 'which source tags went up in tq last month?'"
./agent "which source tags went up in tq last month?"

echo -e "\nTesting 'which partners have been on the yer the most in the last 6 months?'"
./agent "which partners have been on the yer the most in the last 6 months?"

echo -e "\nTesting 'what traffic sources does DNC use?'"
./agent "what traffic sources does DNC use?" 