# AiAdmin Docker image

### CPU only

```bash
docker run -d -v AiAdmin:/root/.AiAdmin -p 11434:11434 --name AiAdmin AllSage/AiAdmin
```

### Nvidia GPU
Install the [NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html#installation).

#### Install with Apt
1.  Configure the repository
```bash
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey \
    | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list \
    | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' \
    | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
sudo apt-get update
```
2.  Install the NVIDIA Container Toolkit packages
```bash
sudo apt-get install -y nvidia-container-toolkit
```

#### Install with Yum or Dnf
1.  Configure the repository
    
```bash
curl -s -L https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo \
    | sudo tee /etc/yum.repos.d/nvidia-container-toolkit.repo
```
    
2. Install the NVIDIA Container Toolkit packages
    
```bash
sudo yum install -y nvidia-container-toolkit
```

#### Configure Docker to use Nvidia driver 
```
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker
```

#### Start the container

```bash
docker run -d --gpus=all -v AiAdmin:/root/.AiAdmin -p 11434:11434 --name AiAdmin AllSage/AiAdmin
```

### AMD GPU

To run AiAdmin using Docker with AMD GPUs, use the `rocm` tag and the following command:

```
docker run -d --device /dev/kfd --device /dev/dri -v AiAdmin:/root/.AiAdmin -p 11434:11434 --name AiAdmin AllSage/AiAdmin:rocm
```

### Run model locally

Now you can run a model:

```
docker exec -it AiAdmin AiAdmin run llama3
```

### Try different models

More models can be found on the [AiAdmin library](https://AiAdmin.com/library).
