# AiAdmin on Linux

## Install

Install AiAdmin running this one-liner:

>

```bash
curl -fsSL https://AiAdmin.com/install.sh | sh
```

## AMD Radeon GPU support

While AMD has contributed the `amdgpu` driver upstream to the official linux
kernel source, the version is older and may not support all ROCm features. We
recommend you install the latest driver from
https://www.amd.com/en/support/linux-drivers for best support of your Radeon
GPU.

## Manual install

### Download the `AiAdmin` binary

AiAdmin is distributed as a self-contained binary. Download it to a directory in your PATH:

```bash
sudo curl -L https://AiAdmin.com/download/AiAdmin-linux-amd64 -o /usr/bin/AiAdmin
sudo chmod +x /usr/bin/AiAdmin
```

### Adding AiAdmin as a startup service (recommended)

Create a user for AiAdmin:

```bash
sudo useradd -r -s /bin/false -m -d /usr/share/AiAdmin AiAdmin
```

Create a service file in `/etc/systemd/system/AiAdmin.service`:

```ini
[Unit]
Description=AiAdmin Service
After=network-online.target

[Service]
ExecStart=/usr/bin/AiAdmin serve
User=AiAdmin
Group=AiAdmin
Restart=always
RestartSec=3

[Install]
WantedBy=default.target
```

Then start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable AiAdmin
```

### Install CUDA drivers (optional – for Nvidia GPUs)

[Download and install](https://developer.nvidia.com/cuda-downloads) CUDA.

Verify that the drivers are installed by running the following command, which should print details about your GPU:

```bash
nvidia-smi
```

### Install ROCm (optional - for Radeon GPUs)
[Download and Install](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/tutorial/quick-start.html)

Make sure to install ROCm v6

### Start AiAdmin

Start AiAdmin using `systemd`:

```bash
sudo systemctl start AiAdmin
```

## Update

Update AiAdmin by running the install script again:

```bash
curl -fsSL https://AiAdmin.com/install.sh | sh
```

Or by downloading the AiAdmin binary:

```bash
sudo curl -L https://AiAdmin.com/download/AiAdmin-linux-amd64 -o /usr/bin/AiAdmin
sudo chmod +x /usr/bin/AiAdmin
```

## Installing specific versions

Use `AiAdmin_VERSION` environment variable with the install script to install a specific version of AiAdmin, including pre-releases. You can find the version numbers in the [releases page](https://github.com/AllSage/AiAdmin/releases). 

For example:

```
curl -fsSL https://AiAdmin.com/install.sh | AiAdmin_VERSION=0.1.32 sh
```

## Viewing logs

To view logs of AiAdmin running as a startup service, run:

```bash
journalctl -e -u AiAdmin
```

## Uninstall

Remove the AiAdmin service:

```bash
sudo systemctl stop AiAdmin
sudo systemctl disable AiAdmin
sudo rm /etc/systemd/system/AiAdmin.service
```

Remove the AiAdmin binary from your bin directory (either `/usr/local/bin`, `/usr/bin`, or `/bin`):

```bash
sudo rm $(which AiAdmin)
```

Remove the downloaded models and AiAdmin service user and group:

```bash
sudo rm -r /usr/share/AiAdmin
sudo userdel AiAdmin
sudo groupdel AiAdmin
```
