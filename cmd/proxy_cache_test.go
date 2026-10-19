package cmd

import (
	"net/http"
	"testing"
)

const (
	testProxyCachePath = "/organization/" + testNamespace + "/proxycache"
	testProxyUpstream  = "docker.io"
	testProxyJSON      = `{"upstream_registry":"docker.io","expiration_s":3600}`
)

func setProxyCacheVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &upstreamRegistry, testProxyUpstream)
	setCommandTestValue(t, &insecure, true)
	setCommandTestValue(t, &expiration, 3600)
	setCommandTestValue(t, &confirm, true)
}

func TestProxyCacheCommands(t *testing.T) {
	tests := []commandDomainCase{
		{
			name: "proxy cache info", run: proxyCacheCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testProxyCachePath},
			response: testProxyJSON, wantOutput: testProxyUpstream, setup: setProxyCacheVars,
		},
		apiErrorCase("proxy cache info", proxyCacheCmd.RunE, testHTTPGet, testProxyCachePath, http.StatusNotFound, "getting proxy cache config", setProxyCacheVars),
		{
			name: "proxy cache create", run: createProxyCacheCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPost, path: testProxyCachePath, body: map[string]interface{}{
				"upstream_registry": testProxyUpstream, "insecure": true, "expiration": 3600,
			}},
			response: testProxyJSON, wantOutput: testProxyUpstream, setup: setProxyCacheVars,
		},
		apiErrorCase("proxy cache create", createProxyCacheCmd.RunE, testHTTPPost, testProxyCachePath, http.StatusBadRequest, "creating proxy cache config", setProxyCacheVars),
		{
			name: "proxy cache delete", run: deleteProxyCacheCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testProxyCachePath, status: http.StatusNoContent},
			setup:   setProxyCacheVars,
		},
		apiErrorCase("proxy cache delete", deleteProxyCacheCmd.RunE, testHTTPDelete, testProxyCachePath, http.StatusNotFound, "deleting proxy cache config", setProxyCacheVars),
	}
	runCommandDomainCases(t, tests)
}

func TestProxyCacheDeleteRequiresConfirm(t *testing.T) {
	setProxyCacheVars(t)
	setCommandTestValue(t, &confirm, false)
	err := runCommandExpectingNoRequest(t, deleteProxyCacheCmd.RunE)
	if err == nil || err.Error() != "must pass --confirm to delete proxy cache config" {
		t.Fatalf("error = %v, want confirm requirement", err)
	}
}

func TestProxyCacheRequiredFlags(t *testing.T) {
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "proxy cache create", args: []string{cmdGet, cmdOrganization, "create-proxy-cache"}, want: []string{testRequiredFlagError, "upstream-registry"}},
	})
}
