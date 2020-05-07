# proximity 👥
### Tracking proximity interactions

### Requirements
- Docker 1.13 (https://www.docker.com/products/docker-desktop)
- Firebase (https://firebase.google.com/)
### Installation
1. Install Docker
2. Pull the Docker image
```
docker pull safedistance/proximity
```
3. Ask Matt to send you the Firebase service account JSON. Save it somewhere accessible
4. Create a "env.list" file with the following variables
### Local execution
Execute the following command to run the service
```
docker run -it -p 127.0.0.1:8080:8080 safedistance/proximity "$(cat path/to/service/account.json)"
```
The service will be accessible at ```localhost:8080/api/v1/proximity```

