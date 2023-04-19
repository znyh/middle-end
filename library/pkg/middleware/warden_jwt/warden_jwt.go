package warden_jwt

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MapClaims map[string]interface{}

type WardenJWT struct {
	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is TokenTime + MaxRefresh.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on login info.
	// Must return user data as user identifier, it will be stored in Claim Array. Required.
	// Check error (e) to determine the appropriate error message.
	Authenticator func(userinfo interface{}) (interface{}, error)

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(data interface{}) MapClaims

	// User can define own RefreshResponse func.
	RefreshResponse func(interface{}, string) (interface{}, error)

	// Set the identity key
	IdentityKey string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time

	// Private key file for asymmetric algorithms
	PrivKeyFile string

	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey
}

var (
	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	ErrMissingAuthenticatorFunc = errors.New("BMJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	ErrMissingLoginValues = errors.New("missing Username or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrMissingExpField missing exp field in token
	ErrMissingExpField = errors.New("missing exp field")

	// ErrWrongFormatOfExp field must be float64 format
	ErrWrongFormatOfExp = errors.New("exp must be float64 format")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = errors.New("parameter token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")

	// IdentityKey default identity key
	IdentityKey = "identity"
)

//func New(wj *WardenJWT) (*WardenJWT, error) {
//	if err := wj.MiddlewareInit(); err != nil {
//		fmt.Println(err)
//		return nil, err
//	}
//
//	return m, nil
//}

func New(wj *WardenJWT) (*WardenJWT, error) {
	if err := wj.WardenJWTInit(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return wj, nil
}

func (wj *WardenJWT) WardenJWTInit() error {

	if wj.SigningAlgorithm == "" {
		wj.SigningAlgorithm = "HS256"
	}

	if wj.Timeout == 0 {
		wj.Timeout = time.Hour
	}

	if wj.TimeFunc == nil {
		wj.TimeFunc = time.Now
	}

	wj.TokenHeadName = strings.TrimSpace(wj.TokenHeadName)
	if len(wj.TokenHeadName) == 0 {
		wj.TokenHeadName = "Bearer"
	}
	//token载荷添加字段IdentityKey 我们的业务里使用userid作为IdentityKey
	//加入后 结构如下{"alg":"RS256","typ":"JWT"}{"exp":1574925240,"orig_iat":1574924940,"userid":"yuki"}
	if wj.IdentityKey == "" {
		wj.IdentityKey = IdentityKey
	}

	if wj.usingPublicKeyAlgo() {
		return wj.readKeys()
	}

	if wj.Key == nil {
		return ErrMissingSecretKey
	}
	return nil
}

func (wj *WardenJWT) MiddlewareFunc() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		//println("=========================GRPC middleware  handler之前执行============================")
		if args.FullMethod != "/hys_auth.v1.Auth/GetToken" {
			error := wj.middlewareImpl(ctx)
			if error != nil {
				return nil, status.Errorf(codes.InvalidArgument,
					error.Error())
			}
		}

		resp, err = handler(ctx, req)

		//println("=========================GRPC middleware handler之后执行============================")
		if args.FullMethod == "/hys_auth.v1.Auth/UpdateToken" {
			resp, err = wj.RefreshHandler(ctx, resp)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument,
					err.Error())
			}
		}

		return
	}
}

func (wj *WardenJWT) middlewareImpl(ctx context.Context) error {
	claims, err := wj.GetClaimsFromJWT(ctx)
	if err != nil {
		return err
	}

	if claims["exp"] == nil {
		return ErrMissingExpField
	}

	if _, ok := claims["exp"].(float64); !ok {
		return ErrWrongFormatOfExp
	}

	if int64(claims["exp"].(float64)) < wj.TimeFunc().Unix() {
		return ErrExpiredToken
	}

	return nil
}

func (wj *WardenJWT) NewToken(userinfo interface{}) (token_str string, err error) {

	if wj.Authenticator == nil {
		return "", errors.New("error Authenticator is nil")
	}

	data, err := wj.Authenticator(userinfo)

	if err != nil {
		return "", errors.New("Authenticator return error")
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(wj.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if wj.PayloadFunc != nil {
		for key, value := range wj.PayloadFunc(data) {
			claims[key] = value
		}
	}

	expire := wj.TimeFunc().Add(wj.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = wj.TimeFunc().Unix()
	tokenString, err := wj.signedString(token)

	if err != nil {
		//mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, c))
		return "", errors.New("error....................")
	}

	return tokenString, nil
	//mw.LoginResponse(c, http.StatusOK, tokenString, expire)
}

func (wj *WardenJWT) signedString(token *jwt.Token) (string, error) {
	wj.readKeys()
	var tokenString string
	var err error
	if wj.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(wj.privKey)
	} else {
		tokenString, err = token.SignedString(wj.Key)
	}
	return tokenString, err
}

func (wj *WardenJWT) usingPublicKeyAlgo() bool {
	switch wj.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

func (wj *WardenJWT) readKeys() error {
	err := wj.privateKey()
	if err != nil {
		return err
	}
	err = wj.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (wj *WardenJWT) privateKey() error {
	keyData, err := ioutil.ReadFile(wj.PrivKeyFile)
	if err != nil {
		return errors.New("public key file unreadable")
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return errors.New("public key invalid")
	}
	wj.privKey = key
	return nil
}

func (wj *WardenJWT) publicKey() error {
	keyData, err := ioutil.ReadFile(wj.PubKeyFile)
	if err != nil {
		return errors.New("public key file unreadable")
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return errors.New("public key invalid")
	}
	wj.pubKey = key
	return nil
}

func (wj *WardenJWT) GetClaimsFromJWT(ctx context.Context) (MapClaims, error) {
	token, err := wj.ParseToken(ctx)

	if err != nil {
		return nil, err
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

func (wj *WardenJWT) ParseToken(ctx context.Context) (*jwt.Token, error) {
	var token string
	var err error

	token, err = wj.jwtFromHeader(ctx, "authorization")

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(wj.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		if wj.usingPublicKeyAlgo() {
			return wj.pubKey, nil
		}

		// save token string if vaild
		//c.Set("JWT_TOKEN", token)

		return wj.Key, nil
	})
}

func (wj *WardenJWT) jwtFromHeader(ctx context.Context, key string) (string, error) {
	var authHeader string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(md)
		if v, ok := md[key]; ok {
			authHeader = v[0]
		}
	}
	//authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == wj.TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (wj *WardenJWT) RefreshHandler(ctx context.Context, resp interface{}) (res interface{}, error error) {
	tokenString, _, err := wj.RefreshToken(ctx)
	if err != nil {
		return nil, err
	}

	res, error = wj.RefreshResponse(resp, tokenString)

	return res, error
}

// RefreshToken refresh token and check if token is expired
func (wj *WardenJWT) RefreshToken(ctx context.Context) (string, time.Time, error) {
	claims, err := wj.CheckIfTokenExpire(ctx)
	if err != nil {
		return "", time.Now(), err
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(wj.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := wj.TimeFunc().Add(wj.Timeout)
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = wj.TimeFunc().Unix()
	tokenString, err := wj.signedString(newToken)

	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expire, nil
}

// CheckIfTokenExpire check if token expire
func (wj *WardenJWT) CheckIfTokenExpire(ctx context.Context) (jwt.MapClaims, error) {
	token, err := wj.ParseToken(ctx)

	if err != nil {
		// If we receive an error, and the error is anything other than a single
		// ValidationErrorExpired, we want to return the error.
		// If the error is just ValidationErrorExpired, we want to continue, as we can still
		// refresh the token if it's within the MaxRefresh time.
		// (see https://github.com/appleboy/gin-jwt/issues/176)
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < wj.TimeFunc().Add(-wj.MaxRefresh).Unix() {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
