# DMX Control

This is a dead simple dmx controller

## Operating Principle

You define a json file (see config.json) of fixtures and commands (ie on,off,etc)  Each state has the channels (1-512) and values (0-255) that should be sent on executing the command. In this way you can create simple on/off commands as well as complex scenes.

The ui will parse the json file and present a table with each fixture and the states as buttons. The controller keeps track of the states and updates the dmx system via usb serial converter when a button is pushed and sends updates periodically to keep things in sync.

![Interface Control](https://github.com/rob121/dmxcontrol/blob/main/readme.png?raw=true)

## Install

Stick the binary (from releases) where you want, the config file will be looked for in the current directory where the program is started or /etc/dmxcontrol
