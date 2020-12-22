# sdup-hue
Converts a phillips Hue API to SDUP

Subscribed Device Update Protocol, SDUP, is a hoppyist project created to consolidate different home automation and IOT control protocols under one umbrella.
The main idea is to convert all exposed APIs and control protocols to one protocol (through proxies called sdup converters) such that all devices, regardless of 
functionality can be treated in a similar fashion.

Every bridge (or device if the device exposes a control protocol) is converted to SDUP and is then controlled by a controller that treats all devices similarly.

Each device is controlled in isolation.
Every device has a set of attributes with a state which in turn can be modified through capabilities.
Attribute examples;

* Active (indicates whether the device is on or off)
* Color (The color the device currently emits)

In turn, related capabilities might be;

* Activate (Changes a boolean attribute to active)
* Deactivate (Changes a boolean state to inactive)
* SetColor (Changes a color attribute to a specified color (requires input))

Worh noting is that the capabilities are attached to an attribute, meaning that multiple attributes may have the capability "activate".

Attributes is something you are while capabilities is something that can be done with that attribute. 
For example; A lamp often has the boolean attribute "active", which indicates whether it is currently shining or not. "active" can either be true or false.
However, "active" can be modified by turning the light on, or turning the light off, "activate" and "deactivate" respectively.

A secondary consideration when developing SDUP is that changes to devices should be pushed to all interested partis, and no poll mechanic should be
implemented. In the case of Phillips Hue, the API does not have any push notification functionality, meaning that SDUP that requires it implemented a poll technique
that is seemlessly translated into a stream of events representing device state changes.
