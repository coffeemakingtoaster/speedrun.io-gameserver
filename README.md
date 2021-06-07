# speedrun.io-gameserver
## Disclaimer
This project is under construction and not finished. Therefore some functionalities might not be implemented to their full extend. This is partly due to the game client (which is on a private gitlab) not being fully finished.
<br> Also as of now the docker build is broken. While it compiles properly there are some issues with the user access rights. This is a known issue that will be fixed in the near future. 
<br> For information on how to deploy the current version of the project see the wiki.

## Overview
This is the server responsible for the real time multiplayer from speedrun.io. 
It is written in golang and based on websockets. This version requires TLS and therefore the corresponding certificates. However it is possible (while not recommended) to deploy this without tls enabled.
For closer information on the codebase and the networking see the wiki.
