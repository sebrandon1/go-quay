package cmd

import (
	"strings"
	"testing"
)

func TestManifestCommands(t *testing.T) {
	const manifestPath = "/repository/testns/testrepo/manifest/sha256:abc"
	baseSetup := func(t *testing.T) {
		setCommandTestValue(t, &namespace, testNamespace)
		setCommandTestValue(t, &repository, testRepository)
		setCommandTestValue(t, &manifestRef, "sha256:abc")
	}
	tests := []commandDomainCase{
		{name: subcmdInfo, run: manifestInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: manifestPath}, response: `{"digest":"sha256:abc"}`, wantOutput: "sha256:abc", setup: baseSetup},
		{name: subcmdDelete, run: manifestDeleteCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: manifestPath}, setup: func(t *testing.T) { baseSetup(t); setCommandTestValue(t, &confirmManifestDeletion, true) }},
		{name: "labels", run: manifestLabelsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: manifestPath + "/labels"}, response: `{"labels":[{"id":"label-1","key":"owner","value":"` + testManifestLabelValue + `"}]}`, wantOutput: testManifestLabelValue, setup: baseSetup},
		{name: "label", run: manifestLabelCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: manifestPath + "/labels/label-1"}, response: `{"id":"label-1","key":"owner","value":"` + testManifestLabelValue + `"}`, wantOutput: "label-1", setup: func(t *testing.T) { baseSetup(t); setCommandTestValue(t, &labelID, "label-1") }},
		{name: "add label", run: manifestAddLabelCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: manifestPath + "/labels", body: map[string]interface{}{"key": "owner", "value": testManifestLabelValue, "media_type": "text/plain"}}, response: `{"id":"label-1","key":"owner","value":"` + testManifestLabelValue + `"}`, wantOutput: testManifestLabelValue, setup: func(t *testing.T) {
			baseSetup(t)
			setCommandTestValue(t, &labelKey, "owner")
			setCommandTestValue(t, &labelValue, testManifestLabelValue)
			setCommandTestValue(t, &labelMediaType, "text/plain")
		}},
		{name: "remove label", run: manifestRemoveLabelCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: manifestPath + "/labels/label-1"}, setup: func(t *testing.T) { baseSetup(t); setCommandTestValue(t, &labelID, "label-1") }},
		{name: testAPIErrorCase, run: manifestInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: manifestPath, status: 500}, response: testAPIErrorResponse, wantError: "getting manifest information", setup: baseSetup},
	}
	runCommandDomainCases(t, tests)
}

func TestManifestDeleteRequiresConfirmation(t *testing.T) {
	setCommandTestValue(t, &namespace, testNamespace)
	setCommandTestValue(t, &repository, testRepository)
	setCommandTestValue(t, &manifestRef, "sha256:abc")
	setCommandTestValue(t, &confirmManifestDeletion, false)
	setCommandTestValue(t, &token, testTokenValue)
	setCommandTestValue(t, &quayURL, testUnreachableURL)
	if err := manifestDeleteCmd.RunE(nil, nil); err == nil || !strings.Contains(err.Error(), "Use --confirm to proceed") {
		t.Fatalf("manifest delete error = %v, want confirmation prompt", err)
	}
}
