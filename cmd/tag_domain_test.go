package cmd

import (
	"net/http"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testTagBasePath = "/repository/" + testNamespace + "/" + testRepository + "/tag"
	testTagPath     = testTagBasePath + "/" + testLatestTagName
	testTagDigest   = "sha256:abc123"
	testTagJSON     = `{"name":"latest","manifest_digest":"sha256:abc123","size":42}`
	testTagExpiry   = "2030-01-01T00:00:00Z"
)

func setTagVars(t *testing.T) {
	setCommandTestValue(t, &namespace, testNamespace)
	setCommandTestValue(t, &repository, testRepository)
	setCommandTestValue(t, &tagName, testLatestTagName)
	setCommandTestValue(t, &tagExpiration, testTagExpiry)
	setCommandTestValue(t, &manifestDigest, testTagDigest)
	setCommandTestValue(t, &confirmTagDeletion, true)
	setCommandTestValue(t, &tagPage, 0)
	setCommandTestValue(t, &tagLimit, 0)
}

func TestTagDomainCommands(t *testing.T) {
	digestBody := map[string]interface{}{"manifest_digest": testTagDigest}
	tests := []commandDomainCase{
		{
			name: "tag list", run: tagListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testTagBasePath + "/"},
			response: `{"tags":[` + testTagJSON + `]}`, wantOutput: testLatestTagName, setup: setTagVars,
		},
		apiErrorCase("tag list", tagListCmd.RunE, testHTTPGet, testTagBasePath+"/", http.StatusNotFound, "listing tags", setTagVars),
		apiErrorCase("tag info", tagInfoCmd.RunE, testHTTPGet, testTagPath, http.StatusNotFound, "getting tag information", setTagVars),
		{
			name: "tag update", run: tagUpdateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testTagPath, body: map[string]interface{}{"expiration": testTagExpiry}},
			response: testTagJSON, wantOutput: testLatestTagName, setup: setTagVars,
		},
		apiErrorCase("tag update", tagUpdateCmd.RunE, testHTTPPut, testTagPath, http.StatusBadRequest, "updating tag", setTagVars),
		{
			name: "tag delete", run: tagDeleteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testTagPath, status: http.StatusNoContent},
			setup:   setTagVars,
		},
		apiErrorCase("tag delete", tagDeleteCmd.RunE, testHTTPDelete, testTagPath, http.StatusNotFound, "deleting tag", setTagVars),
		{
			name: "tag history", run: tagHistoryCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testTagPath + "/history"},
			response: `{"tags":[` + testTagJSON + `]}`, wantOutput: testLatestTagName, setup: setTagVars,
		},
		apiErrorCase("tag history", tagHistoryCmd.RunE, testHTTPGet, testTagPath+"/history", http.StatusNotFound, "getting tag history", setTagVars),
		{
			name: "tag revert", run: tagRevertCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testTagPath + "/revert", body: digestBody},
			response: testTagJSON, wantOutput: testLatestTagName, setup: setTagVars,
		},
		apiErrorCase("tag revert", tagRevertCmd.RunE, testHTTPPost, testTagPath+"/revert", http.StatusBadRequest, "reverting tag", setTagVars),
		apiErrorCase("tag change", tagChangeCmd.RunE, testHTTPPut, testTagPath, http.StatusBadRequest, "changing tag", setTagVars),
		{
			name: "tag restore", run: tagRestoreCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPost, path: testTagPath + "/restore", body: digestBody},
			setup:   setTagVars,
		},
		apiErrorCase("tag restore", tagRestoreCmd.RunE, testHTTPPost, testTagPath+"/restore", http.StatusBadRequest, "restoring tag", setTagVars),
		{
			name:        "tag list table",
			executeArgs: []string{cmdGet, cmdTag, subcmdList, "-n", testNamespace, "-r", testRepository, "--output", outputTable},
			request:     commandRequestExpectation{method: testHTTPGet, path: testTagBasePath + "/"},
			response:    `{"tags":[` + testTagJSON + `]}`, wantOutput: "DIGEST",
		},
	}
	runCommandDomainCases(t, tests)
}

func TestTagDomainRequireTagGuards(t *testing.T) {
	for _, run := range map[string]func(*cobra.Command, []string) error{
		"tag update": tagUpdateCmd.PreRunE, "tag history": tagHistoryCmd.PreRunE,
		"tag revert": tagRevertCmd.PreRunE, "tag restore": tagRestoreCmd.PreRunE,
	} {
		if run == nil {
			continue
		}
		setCommandTestValue(t, &tagName, "")
		if err := run(&cobra.Command{}, nil); err == nil {
			t.Error("expected tag-required error, got nil")
		}
	}
}

func TestTagDomainRequiredFlags(t *testing.T) {
	tag := func(sub string) []string {
		return []string{cmdGet, cmdTag, "-n", testNamespace, "-r", testRepository, "-T", testLatestTagName, sub}
	}
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "tag revert flags", args: tag("revert"), want: []string{testRequiredFlagError, cmdManifest}},
		{name: "tag restore flags", args: tag("restore"), want: []string{testRequiredFlagError, cmdManifest}},
	})
}
