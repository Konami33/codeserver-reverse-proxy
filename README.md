## Steps to run this proxy server:

### 1. Pull the following docker images 

Codeserver-node
```bash
docker pull poridhi/codeserver-node:v1.1
```
Codeserver-python
```bash
docker pull poridhi/codeserver-python:v1.2
```

### 2. Run the code-server containers
```bash
docker run -it -p 7080:8080 poridhi/codeserver-python:v1.2
```

```bash
docker run -it -p 8080:8080 poridhi/codeserver-node:v1.1
```


### 3. Build a docker image of the reverse-proxy-server

Path: r-proxy-np>

```bash
docker build -t r-proxy:v1 .
```
### 4. Run a container based on that image

```bash
docker run -it -p 8000:8000 r-proxy:v1
```

### 5. Check the codeserver

1. Node code-server is running on: http://localhost:8080/?folder=/app

2. Python code-server is running on: http://localhost:7080/?folder=/app

### 6. Check the proxy server for redirection

1. Node code-server: http://localhost:8000/node/?folder=/app 

2. Python code server: http://localhost:8000/python/?folder=/app 