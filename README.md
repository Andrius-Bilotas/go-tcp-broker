# TCP Publisher/Subscriber broker server

Simple broker allowing multiple publishers and subscribers to connect. Routes messages from publishers to all subscribers. 

Port for publisher connections: `8000` \
Port for subscriber connections: `8001` 

## Running

To start the server, run this command: \
`make run` \
This will build the application and start the application from the generated executable \
The executable gets placed inside the `bin/` folder \
To just build the app, run this command: \
`make build`

## Testing

To run the tests, run this command: \
`make test` 
