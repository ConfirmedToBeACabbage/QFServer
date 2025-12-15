- Session one to one encrypted
- Table based crypt only held by handler 
- receiver needs to verbally or in another way get it from the handler the tokens 
- the receiver on a randomized gui with randomized positioning selects (clicks) the tokens 
    - Catch anything recording and set it to dark?
    - Get on a recording session (outlay black) 
- When tokens are selected decryption happens 
- The same math happens, the key is that you need to know the tokens otherwise it wont work

Tokens are all special characters

https://www.freecodecamp.org/news/ascii-table-hex-to-ascii-value-character-code-chart-2/

Something to do with Hex

Technically it's just the selection of tokens. I can randomize the gui also. This can be fun

This can be a linux package! some cli type thing. 

How would I defend what is clicked? In memory?
    - I could randomize the memory space after every click. Aka randomize the dictionary connecting the buttons. Reshuffle in a random way. That way different memory is used

I could use key controls
Use the key as a sender of information to handler
AKA got this in this address space then in this address space and then this address space. All of this
still encrypted in the channel. 

H --> R [H:PROVIDESTOKENLIST] 
R --> C [R:NAVIGATESKEYS]
C --> [C:RANDOMIZEMEMORYSPACEWITHEACHKEYPRESS]
[C:ACTION:CLICK] --> C [C:RANDOMIZEMEMORYSPACE] [C:STORENEXTLOCAL] [C:RANDOMIZEMEMORYSPACE] 
C --> R [C:SENDDICTMEMORY:&MEMORYSPACE]
[C:RANDOMIZEMEMORYSPACE] 

Function list should be a set of pointers. 
All navigation should be basically done in a tree that randomizes its memory space constantly. 

What if I make structs plus functions to return each item. 

Can I have: 
[hidden in memory tightly] <-> [function that accesses but the functions memory space is constantly refreshed]

I could get an illusion of this by creating the type then destroying it. 

[Select] --> [Create & Send]

[Handler should just be hanlding it all] IT's a protocol but in action type thing

1. Handler starts server
2. Receiver enters ip
3. 

Randomize input you get? But the language is only on the handlers side. 

Maybe simple is better for now. Token which then is sent back to the handler and it's sent over. 

Thought it over. 

Sender: 
1. Opens an input stream 
   - Automatically encrypts all input
   - The area where input is given to is allocated memory but is erased and reset every certain amount of keys
   - Or if reading in the text input from a file then it just encrypts on the way in 

1. Sender startup
1.1 Server Start 
1.2 Broadcast on LAN and see if any other devices have this service
1.3 Set them up in gui to send
1.4 Receiver selected
1.4.1 Establish secure ssl line
1.5 Generate private key
(1.6 Encrypt data connection on private key)
1.8 Generate a token and insert into empty key spaces 
(1.8 Use this key to then encrypt the text)
1.9 Read input

2. Sending and decrypt
2.1 Send the encrypted text over the encrypted line to user 
2.2 When user recieves all send user the key with empty key spaces
2.3 Close connection

3. Client decrypt 
3.1 Client would need to know your specific token for that session 
3.2 Enter in token
3.3 Dump into text file
3.4 Remove all data
