package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"gopkg.in/ldap.v2"
)

type APIScheme struct {
	Version string  `json:"apiVersion" binding:"required"`
	Kind    string  `json:"kind" binding:"required"`
	Spec    APISpec `json:"spec" binding:"required"`
}

type APISpec struct {
	Token string `json:"token" binding:"required"`
}

type User struct {
	DN		 string
	Name   string
	ID     string
	Groups []string
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
	var apiScheme APIScheme
	if err := c.ShouldBindJSON(&apiScheme); err == nil {
		user, err := authLDAP(apiScheme.Spec.Token)
		if err == nil && user != nil {
			c.JSON(http.StatusOK, gin.H{
				"apiVersion": "authentication.k8s.io/v1beta1",
				"kind":       "TokenReview",
				"status": gin.H{
					"authenticated": true,
					"user": gin.H{
						"username": user.Name,
						"uid":      user.Name,
						"groups":   user.Groups,
					},
				},
			})
		} else {
			c.JSON(http.StatusUnauthorized, authFailed)
		}
	} else {
		c.JSON(http.StatusUnauthorized, authFailed)
	}
}

func GuidToOctetString(guid string) string {
	var buffer bytes.Buffer
	for index, guidCharacter := range(guid) {
		if index % 2 == 0 {
			buffer.WriteString("\\")
		}
		buffer.WriteString(string(guidCharacter))
	}
	return buffer.String()
}

func authLDAP(token string) (*User, error) {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("[LDAP-AUTH] ")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return nil, err
	}
	defer l.Close()

	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_START_TLS")); ok {
		if err = l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return nil, err
		}
	}

	if err = l.Bind(os.Getenv("BIND_DN"), os.Getenv("BIND_PASSWORD")); err != nil {
		return nil, err
	}

	if strings.Contains(os.Getenv("USER_SEARCH_FILTER"), "objectGUID") {
		token = GuidToOctetString(token)
	}

	sru, err := l.Search(ldap.NewSearchRequest(
		os.Getenv("USER_SEARCH_BASE"), ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(os.Getenv("USER_SEARCH_FILTER"), token),
		[]string{os.Getenv("USER_NAME_ATTRIBUTE"), os.Getenv("USER_UID_ATTRIBUTE")},
		nil,
	))
	if err != nil {
		return nil, err
	}

	// only one user
	if len(sru.Entries) != 1 {
		return nil, fmt.Errorf("too much user response")
	}

	user := &User{}
	for _, entry := range sru.Entries {
		user.Name = entry.GetAttributeValue(os.Getenv("USER_NAME_ATTRIBUTE"))
		user.ID = entry.GetAttributeValue(os.Getenv("USER_UID_ATTRIBUTE"))
		log.Printf("Search user result: %v %v\n", user.Name, user.ID)
	}

	srg, err := l.Search(ldap.NewSearchRequest(
		os.Getenv("GROUP_SEARCH_BASE"), ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(os.Getenv("GROUP_SEARCH_FILTER"), user.Name, user.DN),
		[]string{os.Getenv("GROUP_NAME_ATTRIBUTE")},
		nil,
	))
	if err != nil {
		return nil, err
	}

	for _, entry := range srg.Entries {
		g := entry.GetAttributeValue(os.Getenv("GROUP_NAME_ATTRIBUTE"))
		user.Groups = append(user.Groups, g)
	}
	log.Printf("Search group result: %v\n", user.Groups)
	return user, nil
}

func main() {
	listenAddr := flag.String("listen-addr", ":8087", "Authn service listen address.")
	config := flag.String("config", "ldap-auth.conf", "LDAP auth config file.")
	flag.Parse()

	err := godotenv.Load(*config)
	if err != nil {
		log.Fatal("Error loading config file")
	}

	gin.DisableConsoleColor()
	r := gin.Default()
	r.GET("/healthz", healthz)
	r.POST("/auth", auth)
	r.Run(*listenAddr)
}
