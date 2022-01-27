package api

import (
	"github.com/go-pkgz/auth/token"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v39/github"
	"github.com/joscha-alisch/dyve/internal/core/config"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestUpdateClaimsFunc(t *testing.T) {
	f := getUpdateClaimsFunc(Opts{})

	claims := token.Claims{User: &token.User{Name: "name"}}
	res := f(claims)
	if !cmp.Equal(claims, res) {
		t.Errorf("claims mismatch: %s\n", cmp.Diff(claims, res))
	}
}

func TestUpdateClaimsGroups(t *testing.T) {
	f := getUpdateClaimsFunc(Opts{
		DevConfig: config.DevConfig{
			UseFakeOauth2: true,
			UserGroups:    []string{"extra-group"},
		},
	})

	claims := token.Claims{User: &token.User{Name: "name"}}
	res := f(claims)

	expected := token.Claims{User: &token.User{Name: "name"}}
	expected.User.SetSliceAttr("groups", []string{"extra-group"})
	if !cmp.Equal(expected, res) {
		t.Errorf("claims mismatch: %s\n", cmp.Diff(expected, res))
	}
}

func TestTokenValidator(t *testing.T) {
	f := getTokenValidatorFunc(Opts{})

	claims := token.Claims{User: &token.User{Name: "name"}}
	if f("", claims) != true {
		t.Errorf("user should be allowed")
	}
}

func TestTokenValidatorGitHub(t *testing.T) {
	f := getTokenValidatorFunc(Opts{
		Auth: config.AuthConfig{
			GitHub: config.AuthProviderConfig{
				Enabled: true,
				Org:     "allowed-org",
			},
		},
	})

	claims := token.Claims{User: &token.User{ID: "github_23123", Name: "name", Attributes: map[string]interface{}{
		"orgs": []interface{}{"allowed-org", "other-org"},
	}}}
	if f("", claims) != true {
		t.Errorf("user should be allowed")
	}

	claims = token.Claims{User: &token.User{ID: "github_23123", Name: "name", Attributes: map[string]interface{}{
		"orgs": []interface{}{"other-org"},
	}}}
	if f("", claims) == true {
		t.Errorf("user should not be allowed")
	}
}

func TestSecretFunc(t *testing.T) {
	f := getTokenSecretFunc(Opts{Auth: config.AuthConfig{Secret: "my-secret"}})
	token, err := f("")
	if err != nil {
		t.Fatal("unexpected error")
	}
	if token != "my-secret" {
		t.Fatal("token is wrong")
	}
}

func TestGithubProvider(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.github.com").
		Get("/user/teams").
		Reply(200).
		JSON([]*github.Team{{ID: github.Int64(123), Organization: &github.Organization{
			Login: github.String("organization"),
		}}, {ID: github.Int64(456), Organization: &github.Organization{
			Login: github.String("organization"),
		}}, {ID: github.Int64(789), Organization: &github.Organization{
			Login: github.String("organization2"),
		}}})

	f := getGHProviderFunc()

	res := f(&http.Client{}, token.User{})

	expected := token.User{
		Attributes: map[string]interface{}{
			"orgs": []string{
				"organization",
				"organization2",
			},
			"groups": []string{
				"github:organization2:789",
				"github:organization:123",
				"github:organization:456",
			},
		},
	}

	if !cmp.Equal(expected, res) {
		t.Errorf("result mismatch: %s\n", cmp.Diff(expected, res))
	}
}
