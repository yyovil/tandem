# Instructions
- build the docker image using this cmd
  ```bash
  docker build -f KaliDockerfile -t kali:headless .
  ```
- create and start the container using the following cmd
  ```bash
  docker run -it --privileged --hostname tandem kali:headless bash
  ```
  > NOTE: future releases won't require the privileged option.
- up the **metasploitable3's ub1404** VM (this is for _eval_ purpose) using this cmd
  ```bash
    vagrant up ub1404
  ```
- prepare the env
  ```bash
    cp tandem.example.env tandem.env
    # paste your environment variables here
  ```
  ```bash
    source tandem.env
  ```
- run 
  ```bash 
  tandem
  ```