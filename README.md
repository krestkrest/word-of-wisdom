## Word of wisdom

Simple “Word of Wisdom” TCP-server implementation.  
• TCP is protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol with SHA-1 HashCash is used.  
• After Proof Of Work verification, server sends one of the quotes from the configured file with quotes.  
• Docker files are provided both for the server and for the client that solves the POW challenge

## Build
- make docker-build

## Run
- docker run pow-server
- docker run -it pow-client (pow-client CLI command will be available)