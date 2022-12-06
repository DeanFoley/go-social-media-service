# Dev Diary 5

## Problem

Benchmark tests instantly lock up.

## Pretext

My solution makes use of channels as the spine of concurrent operations; the DB uses a channel-based, synchronous queueing system to ensure that a) instructions are processed in the order of receive, and b) all operations are finished before the server is shut down

## Solution

Waitgroups have not been used in this solution and are *often* used in tandem with channels to ensure synchronous access to resources; this shouldn't be necessary for read-operations, but might be useful when performing write operations (should as establishing a new follow)

Waitgroups might prove useful in ensuring effective concurrency during write-operations, but I haven't had time to go back and refine the implementation. As of now, it appears the system might struuggle to handle the load of too many users creating new accounts and following other people. We would like to politely ask our patrons of Gitter (thats what I'm calling it, because it's like Twitter, but made in Go; aren't I a clever sausage) do not try to do things too fast. Please consult one another to make sure you're not overloading the system. Perhaps establish your own rotas.