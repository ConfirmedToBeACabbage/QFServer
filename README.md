# Abstract

A lan file exchange CLI program. Works with both Windows and Linux. My motivations for this project are simple: 

- Practice Go

The idea came to me when I wanted to pass some notes to my other device on my local network. However, I wanted to do it without exiting the lan network. AKA Gmail or other notes applications. All of it should be in house. 

## Concept

A receiver and sender. Both find each other using broadcasts on the local network, adding it to a pool of known users. When a file is prepared for transfer, you simply choose the corresponding host and send the file over an end to end encrypted exchange. 

## Project Structure


```
QuickFileShareProjectLan
├─ Idea.md
├─ client
│  ├─ broker.go
│  ├─ client.go
│  ├─ command.go
│  └─ go.mod
├─ crypt
│  ├─ crypt.go
│  └─ go.mod
├─ fr
│  ├─ filereader.go
│  └─ go.mod
├─ go.mod
├─ log
│  ├─ go.mod
│  └─ log.go
├─ main.go
├─ readme.md
└─ server
   ├─ go.mod
   └─ server.go

```

## Running & Commands

Once the project is cloned into a local directory, starting it is simple. In the main directory "go run ." should begin the application. 

Running "help" should give you more details on various commands you can use. 

Right now the functionality that works is: 

1. util server open
2. util server broadcast
3. util server pool 

These three items work in sequence to: 

1. Open the server on UDP for listenening and receiving incoming connections 
2. Begin the broadcast process
3. Show the collected addreses in the pool

## TODO

Currently, with the way go-routines are done; Input and output isn't organized. So syntax may
seem all over the place. 

The current TODOs: 

1. Input and Output stream where goroutine output is buffered
2. Debugger and Logger both to be configured with an output stream 
3. Extra commands to flesh out the control of various workers 

## Submitting Changes

I'm currently finishing a semester and won't be able to dedicate much time to this project. However, it's on standby as I take learnings and implement them in other private repositories. 

I won't be able to check pull requests immediately, I do look forward with excitment on any improvements or changes. 

I won't be accepting full refactors of this program, however implementations of extra modules or other smaller changes are welcome. Bonus if you can add a comment explaining a learning. 