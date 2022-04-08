package controllers

import (
	"net/http"
	"os"
	"strconv"

	"k8smanager/internal/clientgo"

	"k8smanager/internal/clientgo/joomla"
	"k8smanager/internal/clientgo/wordpress"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RequestApp struct {
	Name    string `json:"name"`
	Port    string `json:"port"`
	App     string `json:"app"`
	Object  string `json:"object"`
	Replica string `json:"replica"`
}

func AppCreate(c *gin.Context) {
	rac := RequestApp{}
	err := c.ShouldBindJSON(&rac)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		log.Error(err)
		return
	}
	p, err := strconv.Atoi(rac.Port)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ap clientgo.AppInterface
	if rac.App == "wordpress" {
		ap = wordpress.NewWordpressApp()
	} else if rac.App == "joomla" {
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App name"})
		return
	}
	err = ap.Create(rac.Name, int32(p))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(201, gin.H{"Created " + rac.App + " App ": true})
}

func AppDelete(c *gin.Context) {
	rac := RequestApp{}
	err := c.ShouldBindJSON(&rac)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	var ap clientgo.AppInterface
	if rac.App == "wordpress" {
		os.Setenv("NAMESPACE", rac.Name+"-wp")
		ap = wordpress.NewWordpressApp()
	} else if rac.App == "joomla" {
		os.Setenv("NAMESPACE", rac.Name+"-joomla")
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App Name"})
		return
	}
	err = ap.Delete(rac.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(202, gin.H{"Deleted " + rac.App + " App": rac.Name})
}

func ListAll(c *gin.Context) {
	appname := c.Request.URL.Query().Get("app")
	name := c.Request.URL.Query().Get("name")
	var ap clientgo.AppInterface
	if appname == "wordpress" {
		os.Setenv("NAMESPACE", name+"-wp")
		ap = wordpress.NewWordpressApp()
	} else if appname == "joomla" {
		os.Setenv("NAMESPACE", name+"-joomla")
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App Details"})
		return
	}
	servicesName, deploymentsName, podsName, pvcsName, err := ap.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Service List": servicesName,
		"Deployment list": deploymentsName,
		"Pod list":        podsName,
		"PVC List":        pvcsName})
}

func GetDetails(c *gin.Context) {
	appname := c.Request.URL.Query().Get("app")
	object := c.Request.URL.Query().Get("object")
	name := c.Request.URL.Query().Get("name")
	var ap clientgo.AppInterface
	if appname == "wordpress" {
		os.Setenv("NAMESPACE", name+"-wp")
		ap = wordpress.NewWordpressApp()
	} else if appname == "joomla" {
		os.Setenv("NAMESPACE", name+"-joomla")
		ap = wordpress.NewWordpressApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong Details Name"})
		return
	}
	switch object {
	case "deployment":
		Deploymentdetail, _, err := ap.Detail(name, "deployment")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{name + " Details": Deploymentdetail})
		return
	case "service":
		_, Servicedetail, err := ap.Detail(name, "service")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{name + " Details": Servicedetail})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong object Name"})
		return
	}
}

func AppUpdate(c *gin.Context) {
	rac := RequestApp{}
	err := c.ShouldBindJSON(&rac)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	var ap clientgo.AppInterface
	if rac.App == "wordpress" {
		os.Setenv("NAMESPACE", rac.Name+"-wp")
		ap = wordpress.NewWordpressApp()
	} else if rac.App == "joomla" {
		os.Setenv("NAMESPACE", rac.Name+"-joomla")
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App Name"})
	}
	switch rac.Object {
	case "deployment":
		replica, err := strconv.Atoi(rac.Replica)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		err = ap.Update(int32(replica), rac.Name, rac.Object)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"Updated Deployment": "Success"})
		return
	case "service":
		nport, err := strconv.Atoi(rac.Port)
		if err != nil {
			log.Error(err)
			return
		}
		err = ap.Update(int32(nport), rac.Name, rac.Object)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(201, gin.H{"Updated Service Port number": true})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong Object Name"})
	}

}
