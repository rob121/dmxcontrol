# DMX Control

This is a dead simple dmx controller

## Operating Principle

You define a json file (see config.json)_ of fixtures and states (ie on,off,etc)  Each state has the channel and value (0-255) that should be sent on executing the command. In this way you can create simple on/off commands as well as complex scenes.

The ui will parse the json file and present a table with each fixture and the states as buttons. The controller keeps track of the states and updates the dmx system via usb serial converter when a button is pushed and sends updates periodically to keep things in sync.


![Interface Control](https://github.com/rob121/dmxcontrol/blob/maing/readme.png?raw=true)