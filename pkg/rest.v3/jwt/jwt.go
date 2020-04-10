package jwt

import (
	"errors"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"

	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var (
	jwtTimeFormat = "2006-01-02 15:04:05"
	// 默认jwt的过期时间是12小时
	defalutJwtDuration = time.Hour * 12
)

// iris provides some basic middleware, most for your learning courve.
// You can use any net/http compatible middleware with iris.FromStd wrapper.
//
// JWT net/http video tutorial for golang newcomers: https://www.youtube.com/watch?v=dgJFeqeXVKw
//
// Unlike the other middleware, this middleware was cloned from external source: https://github.com/auth0/go-jwt-middleware
// (because it used "context" to define the user but we don't need that so a simple iris.FromStd wouldn't work as expected.)
// jwt_test.go also didn't created by me:
// 28 Jul 2016
// @heralight heralight add jwt unit test.
//
// So if this doesn't works for you just try other net/http compatible middleware and bind it via `iris.FromStd(myHandlerWithNext)`,
// It's here for your learning curve.

// A function called whenever an error is encountered
type errorHandler func(context.Context, string)

// TokenExtractor is a function that takes a context as input and returns
// either a token or an error.  An error should only be returned if an attempt
// to specify a token was found, but the information was somehow incorrectly
// formed.  In the case where a token is simply not present, this should not
// be treated as an error.  An empty string should be returned in that case.
type TokenExtractor func(context.Context) (string, error)

// Middleware the middleware for JSON Web tokens authentication method
type Middleware struct {
	Config Config
}

// OnError default error handler
func OnError(ctx context.Context, err string) {
	ctx.StatusCode(iris.StatusUnauthorized)
	_, errWrite := ctx.Writef(err)
	if errWrite != nil {
		return
	}
}

// New constructs a new Secure instance with supplied options.
func New(cfg ...Config) *Middleware {

	var c Config
	if len(cfg) == 0 {
		c = Config{}
	} else {
		c = cfg[0]
	}

	if c.ContextKey == "" {
		c.ContextKey = DefaultContextKey
	}

	if c.ErrorHandler == nil {
		c.ErrorHandler = OnError
	}

	if c.Extractor == nil {
		c.Extractor = FromAuthHeader
	}

	if c.JwtSignInKey == "" {
		c.JwtSignInKey = "jinmu"
	}
	// JwtDuration没有设置就给默认值
	if c.JwtDuration == 0 {
		c.JwtDuration = defalutJwtDuration
	}
	return &Middleware{Config: c}
}

func (m *Middleware) logf(format string, args ...interface{}) {
	if m.Config.Debug {
		log.Printf(format, args...)
	}
}

// Get returns the user (&token) information for this client/request
func (m *Middleware) Get(ctx context.Context) *jwt.Token {
	return ctx.Values().Get(m.Config.ContextKey).(*jwt.Token)
}

// Serve the middleware's action
func (m *Middleware) Serve(ctx context.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		m.Config.ErrorHandler(ctx, err.Error())
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(ctx context.Context) (string, error) {
	return ctx.GetHeader("Authorization"), nil
}

// FromParameter returns a function that extracts the token from the specified
// query string parameter
func FromParameter(param string) TokenExtractor {
	return func(ctx context.Context) (string, error) {
		return ctx.URLParam(param), nil
	}
}

// FromFirst returns a function that runs multiple token extractors and takes the
// first token it finds
func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(ctx context.Context) (string, error) {
		for _, ex := range extractors {
			token, err := ex(ctx)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}

// BuildJwtToken build jwt token
func (m *Middleware) BuildJwtToken(hmac *jwt.SigningMethodHMAC, clientID, zone, name, customizedCode string) (string, error) {
	expiredAt := time.Now().Add(m.Config.JwtDuration)
	token := jwt.NewWithClaims(hmac, jwt.MapClaims{
		"client_id":       clientID,
		"name":            name,
		"zone":            zone,
		"customized_code": customizedCode,
		"expiredAt":       expiredAt.Format(jwtTimeFormat),
	})
	return token.SignedString([]byte(m.Config.JwtSignInKey))
}

// CheckJWT the main functionality, checks for token
func (m *Middleware) CheckJWT(ctx context.Context) error {
	if !m.Config.EnableAuthOnOptions {
		if ctx.Method() == iris.MethodOptions {
			return nil
		}
	}

	// Use the specified token extractor to extract a token from the request
	token, err := m.Config.Extractor(ctx)

	// If debugging is turned on, log the outcome
	if err != nil {
		m.logf("Error extracting JWT: %v", err)
	} else {
		m.logf("Token extracted: %s", token)
	}

	// If an error occurs, call the error handler and return an error
	if err != nil {
		return fmt.Errorf("Error extracting token: %v", err)
	}

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Config.CredentialsOptional {
			m.logf("  No credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil
		}

		// If we get here, the required token is missing
		errorMsg := "Required authorization token not found"
		m.logf("  Error: No credentials found (CredentialsOptional=false)")
		return fmt.Errorf(errorMsg)
	}

	// Now parse the token

	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil {
		m.logf("Error parsing token: %v", err)
		return fmt.Errorf("Error parsing token: %v", err)
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		m.logf("Error validating token algorithm: %s", message)
		return fmt.Errorf("Error validating token algorithm: %s", message)
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		m.logf("Token is invalid")
		return fmt.Errorf("Token is invalid")
	}

	if m.Config.Expiration {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			expiredAtStr, ok := (claims["expiredAt"]).(string)
			if !ok {
				return errors.New("授权认证失败")
			}
			expiredAt, err := time.Parse(jwtTimeFormat, expiredAtStr)
			if err != nil || expiredAt.Before(time.Now()) {
				return errors.New("授权认证失败")
			}
		}
	}

	m.logf("JWT: %v", parsedToken)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Values().Set(m.Config.ContextKey, parsedToken)

	return nil
}
