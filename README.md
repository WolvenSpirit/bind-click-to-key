## Bind click to key

This is an app that will bind cursor clicks at specified coordinates to a keyboard key.

This is useful if you are extensively using an android emulator on Linux (like genymotion) that doesn't expose by itself a way of mapping the on screen touch buttons (like in games) to keyboard keys.

This provides an alternative way of making it work although it has it's obvious shortcomings which is the fact that there is only one cursor but there may be multiple buttons on screen that need to be simultaneusly pressed, in such a case this tool wont help much.

For convenience this program produces a file named `key_bindings.json` before it exits, which it will use to remember your previous key bindings.

## To people who are still searching for ways to make binding keys work:

- There is also ADB (android debug bridge) that can be used to dispatch touch events. I briefly tried it via shell but latency is greater than using this click to key mapper. I did not find an unix socket API or something similar exposed by ADB. If anybody knows of anything I would appreciate a link and will incorporate it into this small app also to overcome it's current limitations.
