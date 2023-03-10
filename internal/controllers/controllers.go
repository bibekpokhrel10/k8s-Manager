package controllers

import (
	"net/http"
	"regexp"
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
	Enable  string `json:"enable"`
}

func AppCreate(c *gin.Context) {
	reqApp := RequestApp{}
	err := c.ShouldBindJSON(&reqApp)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		log.Error(err)
		return
	}
	re, err := regexp.Compile(`[^(a-z)(A-Z)(0-9)-]`)
	if err != nil {
		log.Fatal(err)
	}
	name := re.ReplaceAllString(reqApp.Name, "")
	p, err := strconv.Atoi(reqApp.Port)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ap clientgo.AppInterface
	if reqApp.App == "wordpress" {
		ap = wordpress.NewWordpressApp()
	} else if reqApp.App == "joomla" {
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App name"})
		return
	}
	err = ap.Create(name, int32(p))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(201, gin.H{"Created " + reqApp.App + " App ": true})
}

func AppDelete(c *gin.Context) {
	reqApp := RequestApp{}
	err := c.ShouldBindJSON(&reqApp)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	var ap clientgo.AppInterface
	if reqApp.App == "wordpress" {
		ap = wordpress.NewWordpressApp()
	} else if reqApp.App == "joomla" {
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App Name"})
		return
	}
	err = ap.Delete(reqApp.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(202, gin.H{"Deleted " + reqApp.App + " App": reqApp.Name})
}

func ListAll(c *gin.Context) {
	appname := c.Request.URL.Query().Get("app")
	name := c.Request.URL.Query().Get("name")
	var ap clientgo.AppInterface
	if appname == "wordpress" {
		ap = wordpress.NewWordpressApp()
	} else if appname == "joomla" {
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App Details"})
		return
	}
	Names, err := ap.List(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Service List": Names.Service,
		"Deployment list": Names.Deployment,
		"Pod list":        Names.Pod,
		"PVC List":        Names.Pvc})
}

func GetDetails(c *gin.Context) {
	appname := c.Request.URL.Query().Get("app")
	name := c.Request.URL.Query().Get("name")
	var ap clientgo.AppInterface
	if appname == "wordpress" {
		ap = wordpress.NewWordpressApp()
	} else if appname == "joomla" {
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong Details Name"})
		return
	}
	Deploymentdetail, Servicedetail, err := ap.Detail(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{name + " Deployment Details ": Deploymentdetail,
		name + " Service Details ": Servicedetail})

}

func AppUpdate(c *gin.Context) {
	reqApp := RequestApp{}
	err := c.ShouldBindJSON(&reqApp)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	var ap clientgo.AppInterface
	if reqApp.App == "wordpress" {
		ap = wordpress.NewWordpressApp()
	} else if reqApp.App == "joomla" {
		ap = joomla.NewJoomlaApp()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong App Name"})
	}
	switch reqApp.Object {
	case "deployment":
		replica, err := strconv.Atoi(reqApp.Replica)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		err = ap.Update(int32(replica), reqApp.Name, reqApp.Object, reqApp.Enable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"Updated Deployment replica to " + reqApp.Replica: "Success"})
		return
	case "service":
		nport, err := strconv.Atoi(reqApp.Port)
		if err != nil {
			log.Error(err)
			return
		}
		err = ap.Update(int32(nport), reqApp.Name, reqApp.Object, reqApp.Enable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(201, gin.H{"Updated Service Port number to " + reqApp.Port: true})
		return
	case "filemanager":
		log.Info("I am inside file top")
		err = ap.Update(int32(1), reqApp.Name, reqApp.Object, reqApp.Enable)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		log.Info("I am inside file")
		c.JSON(201, gin.H{"Status": "Successfull"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong Object Name"})
	}

}
