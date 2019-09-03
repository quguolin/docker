- docker mysql
    ```bash
    docker run --name mysql5.6 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.6
    ```
- docker stop all containers
    ```bash
    docker stop $(docker ps -aq)
    ```  
    
- docker remove all containers
    ```bash
    docker rm $(docker ps -aq)
    ```  
    
- docker remove all images
    ```bash
    docker rmi $(docker images -q)
