package tonoauth

import (
	"context"
	"errors"
	"os"

	"github.com/bindernews/sts-msr/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/oauth2"
)

var ErrInvalidSecret = errors.New("secret cannot be empty")
var ErrBadCsrf = errors.New("CSRF token validation error")
var ErrSessionExpired = errors.New("session expired, please try again")
var ErrBadRequest = errors.New("bad request, please try again")
var ErrAuthFailed = errors.New("oauth failed")
var ErrNoEmail = errors.New("could not get email address")

const (
	keyOauthKind = "tonoAuthKind"
	// Session key for the user's email, which will only be set if they are logged in
	KeyUserEmail = "TonoUserEmail"
	// Session key for user redirect after authentication
	KeyRedirectUrl = "TonoRedirectUrl"
)

type Provider struct {
	// Name of provider, used as the sub-url for redirects
	Name     string
	Config   oauth2.Config
	GetEmail func(context.Context, *oauth2.Token) (string, error)
}

type Controller struct {
	// Default URL to redirect to on login/logout
	HomeUrl string
	// Base route to use for oauth callbacks, must NOT end with a / (default is "/oauth")
	Route string
	// The CSRF secret, default is `os.Getenv("CSRF_SECRET")`
	Secret string
	// Map of provider names to providers
	providers map[string]*Provider
	baseUrl   string
}

func NewOauthController(baseUrl string) *Controller {
	s := Controller{}
	s.baseUrl = baseUrl
	s.providers = make(map[string]*Provider)
	// Default values
	s.Route = "/oauth"
	s.HomeUrl = s.baseUrl
	s.Secret = os.Getenv("CSRF_SECRET")
	return &s
}

// Returns the redirect
func (s *Controller) RedirectUrl() string {
	return s.baseUrl + s.Route + "/redirect"
}

func (s *Controller) AddProviders(providers ...*Provider) *Controller {
	for _, p := range providers {
		p.Config.RedirectURL = s.RedirectUrl()
		s.providers[p.Name] = p
	}
	return s
}

func (s *Controller) Init(r *gin.Engine) error {
	if s.Secret == "" {
		return ErrInvalidSecret
	}
	// Add gin routes
	r.Group(s.Route).
		Use(csrf.Middleware(csrf.Options{
			Secret: s.Secret,
		})).
		GET("/login", s.LoginRedirect).
		GET("/logout", s.LogoutRedirect).
		GET("/redirect", s.GetRedirect)
	return nil
}

func (s *Controller) LoginRedirect(c *gin.Context) {
	sess := sessions.Default(c)
	defer sess.Save()
	var url string

	var params struct {
		Kind string `form:"provider" binding:"required"`
	}

	// Read provider
	if err := c.BindQuery(&params); err != nil {
		util.AbortErr(c, 400, ErrBadRequest)
		return
	}

	// If user is already logged in skip login
	if util.SessGetString(sess, KeyUserEmail) != "" {
		url = s.getNextUrl(sess)
	} else {
		sess.Set(keyOauthKind, params.Kind)
		conf := s.providers[params.Kind].Config
		url = conf.AuthCodeURL(csrf.GetToken(c), oauth2.AccessTypeOffline)
	}
	c.Redirect(303, url)
}

func (s *Controller) LogoutRedirect(c *gin.Context) {
	sess := sessions.Default(c)
	defer sess.Save()
	sess.Delete(KeyUserEmail)
	sess.Delete("csrfToken")
	url := s.getNextUrl(sess)
	c.Redirect(303, url)
}

func (s *Controller) GetRedirect(c *gin.Context) {
	var params struct {
		State string `form:"state" binding:"required"`
		Code  string `form:"code" binding:"required"`
	}
	ctx := c.Request.Context()
	sess := sessions.Default(c)
	// Read params
	if err := c.BindQuery(&params); err != nil {
		util.AbortErr(c, 400, err)
		return
	}
	// csrf check
	if params.State != csrf.GetToken(c) {
		util.AbortErr(c, 403, ErrBadCsrf)
		return
	}
	// Get provider
	provider := s.providers[util.SessGetString(sess, keyOauthKind)]
	if provider == nil {
		util.AbortErr(c, 403, ErrSessionExpired)
		return
	}
	// Obtain the real token
	token, err := provider.Config.Exchange(ctx, params.Code)
	if err != nil || token == nil {
		// Log the error
		c.Error(err)
		util.AbortErr(c, 403, ErrAuthFailed)
		return
	}
	// Get email
	email, err := provider.GetEmail(ctx, token)
	if err != nil {
		util.AbortErr(c, 403, err)
		return
	}

	// Store email in the session, get next URL, and redirect
	sess.Set(KeyUserEmail, email)
	nextUrl := s.getNextUrl(sess)
	defer sess.Save()
	c.Redirect(303, nextUrl)
}

// Get an delete the session-stored redirect URL.
// Make sure to call sess.Save()
func (s *Controller) getNextUrl(sess sessions.Session) string {
	// Get the redirect URL and then clear it
	nextUrl := util.SessGetString(sess, KeyRedirectUrl)
	if nextUrl == "" {
		nextUrl = s.HomeUrl
	}
	sess.Delete(KeyRedirectUrl)
	return nextUrl
}

type ProviderOpts struct {
	ClientID     string
	ClientSecret string
}

func NewProviderOptsFromEnv(prefix string) ProviderOpts {
	return ProviderOpts{
		ClientID:     os.Getenv(prefix + "CLIENT_ID"),
		ClientSecret: os.Getenv(prefix + "CLIENT_SECRET"),
	}
}
