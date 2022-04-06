package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"wordpress.com/internal/clientgo/wordpress"
)

func AppCreate(c *gin.Context) {
	appname := c.Request.URL.Query().Get("app")
	wname := c.Request.URL.Query().Get("name")
	port := c.Request.URL.Query().Get("port")
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Error(err)
		return
	}
	if appname == "wordpress" {
		wp := wordpress.NewWordpressApp(wname)
		err = wp.Create(wname, int32(p))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Created Wordpress App": true})
	}
}

func AppDelete(c *gin.Context) {
	appname := c.Request.URL.Query().Get("app")
	wname := c.Request.URL.Query().Get("name")
	if appname == "wordpress" {
		wp := wordpress.NewWordpressApp(wname)
		err := wp.Delete(wname)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Delete wordpress App": wname})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": "Wrong App Name"})
	}
}

func ListAll(c *gin.Context) {
	wp := wordpress.NewWordpressApp("wordpress")
	servicesName, deploymentsName, podsName, pvcsName := wp.List()

	c.JSON(http.StatusOK, gin.H{"Service List": servicesName,
		"Deployment list": deploymentsName,
		"Pod list":        podsName,
		"PVC List":        pvcsName})
}

func GetDetails(c *gin.Context) {
	oname := c.Request.URL.Query().Get("object")
	dname := c.Request.URL.Query().Get("name")

	wp := wordpress.NewWordpressApp(dname)
	if oname == "deployment" {
		Deploymentdetail, _, err := wp.Detail(dname, "deployment")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{dname + " Details": Deploymentdetail})
	} else if oname == "service" {
		_, Servicedetail, err := wp.Detail(dname, "service")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{dname + " Details": Servicedetail})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": "Wrong object Name"})
	}
}

func AppUpdate(c *gin.Context) {
	oname := c.Request.URL.Query().Get("object")
	switch oname {
	case "deployment":
		dname := c.Request.URL.Query().Get("name")
		replicas := c.Request.URL.Query().Get("replicas")
		replica, err := strconv.Atoi(replicas)
		if err != nil {
			log.Error(err)
			return
		}
		wp := wordpress.NewWordpressApp(dname)
		err = wp.Update(int32(replica), dname, oname)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Updated Deployment": "Success"})
		return
	case "service":
		dname := c.Request.URL.Query().Get("name")
		port := c.Request.URL.Query().Get("port")
		nport, err := strconv.Atoi(port)
		if err != nil {
			log.Error(err)
			return
		}
		wp := wordpress.NewWordpressApp(dname)
		err = wp.Update(int32(nport), dname, oname)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Updated Service Port number": true})
		return
	default:
		c.JSON(http.StatusOK, gin.H{"Error": "Wrong Object Name"})
	}

}
