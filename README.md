# JSON Parser 

## Coding Challenge #2 - Json parsing tool 

This repository holds the source code completing the second challenge from [CodingChallenges.fyi](https://codingchallenges.fyi/challenges/challenge-json-parser).

This tool either successfully parses valid JSON with exit code 0, or fails on invalid JSON with exit code 1. 
The error messages are not too informative yet, but it follows all the rules of you can find on [JSON.org](https://json.org). 

## Usage 

Run the following to make an executable file on MAC/Linux environments: 

`make build-mac`

OR 

Run the following for Windows: 

`make build-windows`

You can run the tool using: 

```bash
./json-parser <filename>
```

OR 

```bash
cat <filename> | ./json-parser
```

## Test

To run the tests, you can use the command: 

`make test`