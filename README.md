# How to start the system
#### Step 1: the log 
Nodes output log messages to log.log in the files directory. This file does not clear itself, it has to be manually deleted, unfortunately :( 

Make sure to delete this file before running any nodes for a clean output. 
#### Step 2: defining the nodes
Start by navigating to the nodes.txt file with `cd src/files`. 
This is where you define the nodes that you want to run, in the format of `IP:portnumber`. 
It is important that you make a go routine for each of the nodes specified in this file.

Now that you have defined them, you should run each of them.

#### Step 3: running the nodes
Open a terminal for each node you want to run. Navigate to src/node with `cd src/node` and run the command `go run node.go`. The terminal will then prompt you to enter a port number. Ensure you enter a valid port number from the nodes.txt file. 

Once every node is run, you can try to access the critical section by typing `cs` in the terminal of a node. This will start the process of the node sending out requests. Try gaining access on multiple nodes.





