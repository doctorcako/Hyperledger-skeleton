package aws

import (
	"context"
	"github.com/lestrrat-go/jwx/jwk"
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
