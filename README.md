# k8s Manager Application

This is a k8s Manager Application.

- Language: Golang

## Requirement

### Golang Requirement

1. Golang https://go.dev/dl/
2. Gin package for router ```go get github.com/gin-gonic/gin ```
3. Godotenv package for env variable ```go get github.com/joho/godotenv ```
4. For logging ``` go get github.com/sirupsen/logrus```
5. For using client-go package 
    - ```go get k8s.io/api```
    - ```go get k8s.io/apimachinery```
	- ```go get k8s.io/client-go```

### Kubernetes Cluster requirement

Create a kubernetes cluster in which the the k8s manager application will run.

## Start the k8s-manager

Form k8s-manager directory, open terminal and run ```go run main.go```

## Walk through the k8s-manager

Open Postman or other such kind of app.

The API will run on ```localhost:4000```

### Creating Application

Enter ```localhost:4000/create``` and in body give the application and name of the app along with port number.

![Create App](/images/createapp.png)

### Listing the pods, pvcs, deployments, services

Enter  ```localhost:4000/list&app=appname&name=nameofapp``` to view the pods, pvcs, deployments and services of the app.

![View List](/images/list.png)

### Updating the Deployment of the Application

Enter ```localhost:4000/update``` and in body give the app, name of the app, object which is deployment and number of replicas you want to create.

![update deployment](/images/updatedeployment.png)

After Updating..

![After Updating](/images/updatelist.png)

### Updating the Service of the Application

Enter ```localhost:4000/update``` and in body give the app, name of the app, object which is service and port number you want to update to.

![update service](/images/serviceupdate.png)

### Deleting the Application

Enter ```localhost:4000/delete``` and in body give the app and name of app to be deleted.

![Deleted Application](/images/delete.png)

After Deleting..

![After Deleting](/images/deletelist.png)

### Testing the Created Application after updating

Enter ```http://172.19.0.5:32222``` in the browser.
The created Application will open. You may have different IP, for that look into the control-panel-node description and search for Internal IP.

![Opening the App in browser](/images/joomlasample.png)
