# ElectroMart ⚡️🛒

## Description 📜
Electro mart is a web application that allows users to buy electronics.
The application is built as a part of the course IDATG2204 Datamodellering og databasesystemer (2024 VÅR) at NTNU.

## Local Development 🛠
1. Install GO from https://golang.org/dl/

2. Install MySQL:
    - Linux: https://www.geeksforgeeks.org/how-to-install-mysql-on-linux/
    - MacOS: https://www.geeksforgeeks.org/how-to-install-mysql-on-macos/
    - Windows: https://www.geeksforgeeks.org/how-to-install-mysql-in-windows/

3. Login to MySQL:
    - Use the following command:
    ```bash
    mysql -u root -p
    ```
    - Enter your password
    - Create a database with the following command:
    ```sql
    CREATE DATABASE databasename;
    ```
    - Create a user with the following command:
    ```sql
    CREATE USER 'username'@'localhost' IDENTIFIED BY 'password';
    ```
    - Grant privileges to the user with the following command:
    ```sql
    GRANT ALL PRIVILEGES ON databasename
    TO 'username'@'localhost';
    ```
    - Create the tables by running the 2204Database.sql script in the database folder

4. Clone the repository:
    - Navigate to the directory where you want to clone the repository
    - Use the following command:
    ```bash
    git clone https://github.com/EskilAl/database_2024
    ```
   
5. Navigate to the project directory:
    ```bash
    cd database_2024
    ```
   
6. Build the project:
    ```bash
    go build
    ```
   
7. Run the project:
    ```bash
    ./database_2024
    ```

## Usage 🤓


## Contact 📧

- [Eskil Alstad](mailto:eskil.alstad@ntnu.no)
- [Erik Bjørnsen](mailto:erbj@stud.ntnu.no)
- [Simon Husås Houmb](mailto:simon.h.houmb@ntnu.no)
- [Maja Melby](mailto:maja.melby@ntnu.no)
- [Jon André Solberg](mailto:jon.a.h.solberg@ntnu.no)
