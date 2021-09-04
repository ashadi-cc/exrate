## INSTALLATION
- copy .env-example to .env in root folder 
- modify API_KEY value with your access key from fixer.io
- install docker to your machine

## RUN APPLICATION
- run docker-compose up -d 
- api service will available at http://localhost:8080

## API ENDPOINT LIST
- GET /project?name=project-name to retrieve new access key by given project name
- GET /convert?access_key=your-key&from=usd&to=eur&value=1 to convert currency 
