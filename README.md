Project Spothify

Author: Aaron Hsu

## Introduction

The purpose of Project Spothify is to recreate the music application Spotify.
Spotify is made to provide a music streaming device as a service that can handle 
high scalability and low latency (fast) music streaming.

It requires a lot of distributed systems concepts including the following...
 * ACID Client-Server Operations
 * P2P Networking
 * Caching popular requests in between servers
 * Distributed Consistent Hashing
 * Distributed Storage Replication
 * Partition Fault Tolerance & Availability
 * and hopefully, good programming style!

### Language: GO
http://golang.org/doc/
It's a clean and efficient programming language that has built in concurrency mechanisms
that allow it to easily write programs that use multicore/networked machines

### Why GO? 
It's the language I used for the Fall 2012 15-440 Distributed Systems at Carnegie Mellon University
Continue practice an interesting language

### Other README's
Consult README_Overview.txt to look at how source & test files are organized
Consult README_TODO.txt to look at what's been accomplished what hasn't