# RPS
Rock paper scissor spock lizard implementation

# How to
Build with `go build` and run with `./rps`

In addition to the game endpoint one can get the current score with `curl -XGET 'http://localhost:4567/score'`

# Notes
Implementation creates a map that for each possible playerhand has a list of what computerHands it will win against. 
Current solution is dependant on the order of the acceptedHands list.



