# Running AiAdmin on NVIDIA Jetson Devices

AiAdmin runs well on [NVIDIA Jetson Devices](https://www.nvidia.com/en-us/autonomous-machines/embedded-systems/) and should run out of the box with the standard installation instructions. 

The following has been tested on [JetPack 5.1.2](https://developer.nvidia.com/embedded/jetpack), but should also work on JetPack 6.0.

- Install AiAdmin via standard Linux command (ignore the 404 error): `curl https://AiAdmin.com/install.sh | sh`
- Pull the model you want to use (e.g. mistral): `AiAdmin pull mistral`
- Start an interactive session: `AiAdmin run mistral`

And that's it!

# Running AiAdmin in Docker

When running GPU accelerated applications in Docker, it is highly recommended to use [dusty-nv jetson-containers repo](https://github.com/dusty-nv/jetson-containers).