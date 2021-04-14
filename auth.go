package main

import (
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ldap/ldap/v3"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Result struct {
	Dn string `json:"dn"`
	Id string `json:"id"`
	Name string `json:"name"`
	Groups []string `json:"groups"`
}


var auth = func(username, password string, c echo.Context) (bool, error) {
	if username == "" {
		return false, nil
	}
	token, err := generateToken(username)
	if err != nil {
		return false, err
	}
	conn, user, err := authLDAP(token)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	conn.SetTimeout(5 * time.Second)
	if err := conn.Bind(user.Dn, password); err != nil {
		log.Printf("LDAP Failed to auth. %s", err)
		return false, err
	} else {
		log.Println("Authenticated success!")
	}

	c.String(http.StatusOK, token)
	return true, nil
}

func authLDAP(token string) (*ldap.Conn, *Result, error) {
	//fmt.Println(token)
	//username := "qianghu"

	decodedToken, _ := jwt.Parse(token, nil)
	claims, _ := decodedToken.Claims.(jwt.MapClaims)
	username := claims["user"]

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return nil, nil, err
	}
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_START_TLS")); ok {
		if err = l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return nil, nil, err
		}
	}
	err = l.Bind(os.Getenv("BIND_DN"), os.Getenv("BIND_PASSWORD"))
	if err != nil {
		return nil, nil, err
	}
	//Search user
	sru, err := l.Search(ldap.NewSearchRequest(
		os.Getenv("BASE_DN"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(os.Getenv("USER_SEARCH_FILTER"), username),
		//[]string{"dn"},
		[]string{os.Getenv("USER_NAME_ATTRIBUTE"), os.Getenv("USER_UID_ATTRIBUTE")},
		nil,
	))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user. %s", err)
	}
	if len(sru.Entries) != 1 {
		return nil, nil, fmt.Errorf("user don't exist. %s", err)
	}
	r := &Result{} 
	for _, entry := range sru.Entries {
		r.Dn = entry.DN
		r.Name = entry.GetAttributeValue(os.Getenv("USER_NAME_ATTRIBUTE"))
		r.Id = entry.GetAttributeValue(os.Getenv("USER_UID_ATTRIBUTE"))
		log.Printf("Search user result: user (%v) uid (%v) \n", r.Name, r.Id)
	}

	//Search group
	srg, err := l.Search(ldap.NewSearchRequest(
		os.Getenv("BASE_DN"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(os.Getenv("GROUP_SEARCH_FILTER"), r.Name),
		[]string{os.Getenv("GROUP_NAME_ATTRIBUTE")},
		nil,
	))
	if err != nil {
		return nil, nil, err
	}

	for _, entry := range srg.Entries {
		r.Groups = entry.GetEqualFoldAttributeValues(os.Getenv("GROUP_NAME_ATTRIBUTE"))
		//g := entry.GetAttributeValue(os.Getenv("GROUP_NAME_ATTRIBUTE"))
		//r.Groups = append(r.Groups, g)
		log.Printf("Search group result: %v\n", r.Groups)
	}

	return l, r, nil
}

func connLDAP() (*ldap.Conn, error) {
	//ref: https://gist.github.com/tboerger/4840e1b5464fc26fbb165b168be23345
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("LDAP_HOST"), os.Getenv("LDAP_PORT")))
	if err != nil {
		return nil, err
	}
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_START_TLS")); ok {
		if err = l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return nil, err
		}
	}
	err = l.Bind(os.Getenv("BIND_DN"), os.Getenv("BIND_PASSWORD"))
	if err != nil {
		return nil, err
	}
	return l, err
}

type Claims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func generateToken(username string) (string, error) {
	expiresTime, _ := time.ParseDuration(os.Getenv("TOKEN_EXPIRES_TIME"))
	claims := Claims{
		User: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresTime).Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: username,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(os.Getenv("APP_KEY")))
	return token, err
}