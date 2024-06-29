# Vulnerable application with Golang

## Technology stack
- NGINX with Mod Security CRS
- Vulnerable Application
- MySQL Database
- Docker Engine
- Docker Compose
- Golang (programming language of vulnerable application)
- Quay.io

## Deployment

### Requirements

#### Depedencies

- Docker Engine
- Docker Compose

#### Port

- 8080, NGINX with Mod Security CRS
- 9090, Vulnerable application
- 3306, MySQL database

### Steps

#### Build image

The code is already integrated into Quay.io to automatically build the image from the source code in GitHub. See the Quay.io registry of this repository at [this link](https://quay.io/repository/jesayafn/vuln-web). If you want to build with your hand, you can build with the command: `docker builds -f Containerfile -t web-vuln:1 .` Do not forget to change the image tag on `docker-compose.yaml` to `web-vuln:1` for the `webapp` service.

#### Deploy the application stack

After installing required the dependencies on your machine, you can simply run the following command: `docker-compose up -d` to up and run the application. After the application is deployed, you can access the application with `http://[Machine-IP]:8080` for accessing the application through NGINX with Mod Security CRS and `http://[Machine-IP]:9090` for accessing the application without any reverse proxy.

## SQL Injection

### Vulnerability and Mitigation of SQL Injection within the Application
This application has two paths, `/vuln-path` and `/secured-path`. `/vuln-path` is written with SQL injection vulnerability and `/secured-path` is written with mitigated from SQL injection.

```go
router.GET("/vuln-path", func(c *gin.Context) {

		db := dbConn()
		defer db.Close()
		username := c.Query("username")
		query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var username string
			var password string
			err := rows.Scan(&id, &username, &password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ID":       id,
				"Username": username,
				"Password": password,
			})
		}
	})
```

In the code snippet above, the code has a vulnerability to SQL injection because the `username` variable that has a value from the `username` parameter of `/vuln-path` path is directly concatenated into an SQL query. You can inject a query like `' OR '1'='1` (tested) to the parameter of username. To prove this vulnerability, you can access the application with `http://[Machine-IP]:9090/vuln-path?username=alice%27%20OR%20%271%27=%271` that inject `' OR '1'='1` on `username` parameter within `/vuln-path` and you will get `200 OK` results with response body with whole data on the database.

```go
router.GET("/secured-path", func(c *gin.Context) {
		db := dbConn()
		defer db.Close()

		username := c.Query("username")

		// Use parameterized, and prepared query to avoid SQL injection
		stmt, err := db.Prepare("SELECT * FROM users WHERE username = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		rows, err := stmt.Query(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var username string
			var password string
			if err := rows.Scan(&id, &username, &password); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ID":       id,
				"Username": username,
				"Password": password,
			})
		}
	})
```

In the code snippet above, the query statement has been pre-compiled and parameterized by default of the prepared statements feature. It means the path `/secured-path` has been mitigated by treating any input from the `username` variable that has a value from the `username` parameter as a parameter for query and prepared the statement for query, not directly concatenated into an SQL query the value of `username` variable. To prove this vulnerability, you can access the application with `http://[Machine-IP]:9090/secured-path?username=alice%27%20OR%20%271%27=%271` that inject `' OR '1'='1` on `username` parameter within `/vuln-path` and you will get `200 OK` results with empty response body.

### Protection with Mod Security CRS

On this deployment of docker compose, I add NGINX with Mod Security CRS to protect the application against any form of security attack including SQL injection. To prove this protection, you can access the application with `http://[Machine-IP]:8080/vuln-path?username=alice%27%20OR%20%271%27=%271` that inject `' OR '1'='1` on `username` parameter within `/vuln-path` and you will get `403 Forbidden` results on browser or cURL. This add protection for the application if the application accidentally creates a mistake in writing code with SQL injection vulnerability.