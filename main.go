package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIScheme struct {
	Version string  `json:"apiVersion" binding:"required"`
	Kind    string  `json:"kind" binding:"required"`
	Spec    APISpec `json:"spec" binding:"required"`
}

type APISpec struct {
	Token string `json:"token" binding:"required"`
}

var authFailed = gin.H{
	"apiVersion": "authentication.k8s.io/v1beta1",
	"kind":       "TokenReview",
	"status": gin.H{
		"authenticated": false,
	},
}

func healthz(c *gin.Context) {
	c.String(200, "ok")
}

func auth(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var apiScheme APIScheme

	if err := c.ShouldBindJSON(&apiScheme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "ldap.k8s.com", 389))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer l.Close()

	if apiScheme.Spec.Token != "test" {
		c.JSON(http.StatusUnauthorized, authFailed)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status": gin.H{
				"authenticated": true,
				"user": gin.H{
					"username": "test",
					"uid":      "test",
					"groups":   "kubernetes-admin",
				},
			},
		})
	}
}

func main() {
	r := gin.Default()
	r.GET("/healthz", healthz)
	r.POST("/auth", auth)
	r.Run(":8087")
}
