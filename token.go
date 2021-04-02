package main

import (
	"github.com/labstack/echo"
	"net/http"
)

type H map[string]interface{}

var authFailed = H{
	"apiVersion": "authentication.k8s.io/v1beta1",
	"kind": "TokenReview",
	"status": H{
		"authenticated": false,
	},
}

type APIScheme struct {
	Version string `json:"apiVersion"`
	Kind string `json:"kind"`
	Spec APISpec `json:"spec"`
}
type APISpec struct {
	Token string `json:"token"`
}

func token(c echo.Context) error {
	var apiScheme APIScheme
	if err := c.Bind(&apiScheme); err == nil {
		_, user, err := authLDAP(apiScheme.Spec.Token)
		if err == nil && user != nil {
			c.JSON(http.StatusOK, H{
				"apiVersion": "authentication.k8s.io/v1beta1",
				"kind":       "TokenReview",
				"status": H{
					"authenticated": true,
					"user": H{
						"username": user.Name,
						"uid":      user.Id,
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
	return nil
}
