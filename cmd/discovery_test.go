package cmd

import (
	"net/http"
	"testing"
)

const (
	testDiscoveryPath    = "/discovery"
	testCapabilitiesPath = "/registry/capabilities"
	testAppInfoPath      = "/app/client-1"
	testEntitiesPath     = "/entities/ali"
)

func setDiscoveryVars(t *testing.T) {
	setCommandTestValue(t, &discoveryClientID, "client-1")
	setCommandTestValue(t, &entityPrefix, "ali")
	setCommandTestValue(t, &includeOrgs, true)
	setCommandTestValue(t, &includeTeams, true)
}

func TestDiscoveryCommands(t *testing.T) {
	tests := []commandDomainCase{
		{
			name: "discovery api", run: discoveryAPICmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testDiscoveryPath},
			response: `{"version":"1.0","apis":{}}`, setup: setDiscoveryVars,
		},
		apiErrorCase("discovery api", discoveryAPICmd.RunE, testHTTPGet, testDiscoveryPath, http.StatusInternalServerError, "getting discovery", setDiscoveryVars),
		{
			name: "discovery capabilities", run: discoveryCapabilitiesCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testCapabilitiesPath},
			response: `{"sparse_manifests":{"supported":true}}`, setup: setDiscoveryVars,
		},
		apiErrorCase("discovery capabilities", discoveryCapabilitiesCmd.RunE, testHTTPGet, testCapabilitiesPath, http.StatusInternalServerError, "getting registry capabilities", setDiscoveryVars),
		{
			name: "discovery app info", run: discoveryAppInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testAppInfoPath},
			response: `{"name":"ci-app","organization":"testns"}`, wantOutput: "ci-app", setup: setDiscoveryVars,
		},
		apiErrorCase("discovery app info", discoveryAppInfoCmd.RunE, testHTTPGet, testAppInfoPath, http.StatusNotFound, "getting app info", setDiscoveryVars),
		{
			name: "discovery entities", run: discoveryEntitiesCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testEntitiesPath, query: map[string]string{"includeOrgs": testTrueValue, "includeTeams": testTrueValue}},
			response: `{"results":[{"name":"` + testPermUser + `","kind":"user"}]}`, wantOutput: testPermUser, setup: setDiscoveryVars,
		},
		apiErrorCase("discovery entities", discoveryEntitiesCmd.RunE, testHTTPGet, testEntitiesPath, http.StatusBadRequest, "getting entities", setDiscoveryVars),
	}
	runCommandDomainCases(t, tests)
}

func TestDiscoveryRequiredFlags(t *testing.T) {
	disc := func(sub string) []string { return []string{cmdGet, "discovery", sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "discovery app info", args: disc("app-info"), want: []string{testRequiredFlagError, testFlagClientID}},
		{name: "discovery entities", args: disc("entities"), want: []string{testRequiredFlagError, "prefix"}},
	})
}
