# Preconditions 
## Docker 
You need docker with version 4.0 or higher. You can download and install docker from [here](https://www.docker.com/products/docker-desktop)

You can check your docker version with the following command 

    docker --version
## Docker compose
To start the database container you need docker compose. You can check if docker compose is installed with the following command 

    docker-compose --version

Usually you do not need to install docker compose seperatly. 

# Run databases
To run the databases we use docker and docker compose. 

## Start databases 
To start the databases open a cmd or shell wihthin this folder and run the following command. 

    docker-compose up

First start take some time since docker images for databases must be downloaded first. 

### MSSql
MSSql is host on http://localhost:1433. There is no UI to interact with the database. You need an external tool to connect with the database. 

ToDO -> Add description for external tool
### Neo4J
Neo4J database is hosted on http://localhost:7474/. Since there is an build in WebUI you can access the database directly with your browser.  

*Note*: For first start up you will be asked for a password. The initial password is `neo4j` and you have to change it to login for the first time.
### MongoDB
MongoDB datbase is hosted on http://localhost:27017. However, you can not access it directly with your browser. User and password are `admin` and `admin`. 
To give a greater experience we added mongo express to provide a UI to interact with the database. Mongo express can be accessed via http://localhost:8081/ 

# Shutdown databases 
To shutdown running databases you can press `Strg+C` wihtin the cmd or shell you start up them. 
You can also execute `docker-compose down` to shutdown databases.