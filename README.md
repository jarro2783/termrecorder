Distributed terminal recorder client and server.

Record a terminal session or watch another session live. The client sends a
session to the server; the server records the session and makes it available
for viewing by others.

It is possible to directly watch a session from the command line using the
client, but you probably don't want to try to record a session directly.
The client is intended to be started by another process and is sent data
through standard input.
