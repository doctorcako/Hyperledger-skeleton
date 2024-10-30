package aws

import (
	"context"
	"github.com/lestrrat-go/jwx/jwk"
	"repo.plexus.services/1329-004_incibe_reto06/utils/golang/customError"
)

func GetAWSCognitoJwk(region, userPoolId string) (jwk.Set, customError.Error) {
	jwkURL := "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId + "/.well-known/jwks.json"
	awsJwk, err := jwk.Fetch(context.TODO(), jwkURL)
	if err != nil {
		return awsJwk, customError.NewError(customError.AuthorisationError, err.Error())
	}

	return awsJwk, nil
}

func GetAWSCognitoJwkFromFile(path string) (jwk.Set, customError.Error) {
	var awsJwk jwk.Set
	awsJwk, err := jwk.ReadFile(path)
	if err != nil {
		return awsJwk, customError.NewError(customError.AuthorisationError, err.Error())
	}

	return awsJwk, nil
}
