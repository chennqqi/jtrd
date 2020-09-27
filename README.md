jtrd
=============

This repository contains a **Dockerfile** of [hmbd](https://github.com/chennqqi/jtrd/) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/u/sort/jtrd/) published to the public [DockerHub](https://hub.docker.com/).

### Dependencies

-	[kalilinux/kali-linux-docker](https://hub.docker.com/r/kalilinux/kali-linux-docker/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [hmbd](https://github.com/chennqqi/jtrd/) for [Docker](https://www.docker.com/)'s [trusted build](https://index.docker.io/u/sort/jtrd/) published to the public [DockerHub](https://hub.docker.com/).

### Usage

build

	install golang(version>1.9)
	git clone https://github.com/chennqqi/jtrd.git -b docker
	docker build -t jtrd .

run as webservice

	docker run -d -p 8080:8080 sort/jtrd:v1.0.0 web

	curl -F 'filename=@shadow' localhost:8080/file?callback=http://api.xxx.com/result

`timeout` set scan max timeout

`callback` set result call back

	if this param is not set, http call will not return until jtr run finished.
	if this param is set, http call will return immeditally
	Password crack will cost some time, strong proposal get data by callback 
	
run directly	

	docker run --privileged -v$YOURFILEPATH:/tmp/shadow sort/jtrd:v1.0.0 crack /tmp/shadow
	
## Sample Output

### JSON:


```json
	{
	  "tid": "8fdcb8b7-b950-44f6-ad4a-2c27efa421c7",
	  "status": 200,
	  "list": [
	    {
	      "user": "testjhon",
	      "pass": "testjhon",
	      "crypt": "$6$YD8vXqbc$9J2RoH5DVGqvC0GZ4kwH3lWpZmWtxsqW/gEQU6VVPsRyHm.fU6pSDvcSzV3MDQ54vYDCoI6vEWg3lrUXee7F90"
	    },
	    {
	      "user": "connector",
	      "pass": "12345678",
	      "crypt": "$6$KD9o3B6U$Ine4pLGNsZOFgrjDuCiUq2FoQN3mGWm7W5p.nmf5M71yPBnV2zq8zElUgwgNASRF0aLpLWo5EJHDFyKEjF1Qr0"
	    }
	  ],
	  "message": "OK"
	}

```

