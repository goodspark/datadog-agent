// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package containers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		actualTags     []string
		expectedTags   []*regexp.Regexp
		expectedOutput string
	}{
		{
			name: "All good",
			actualTags: []string{
				"foo:foo",
				"bar:bar",
				"baz:baz",
			},
			expectedTags: []*regexp.Regexp{
				regexp.MustCompile(`^foo:foo$`),
				regexp.MustCompile(`^bar:`),
				regexp.MustCompile(`:baz$`),
			},
			expectedOutput: "",
		},
		{
			name: "Unexpected tag",
			actualTags: []string{
				"foo:foo",
				"bar:bar",
				"baz:baz",
				"qux:qux",
			},
			expectedTags: []*regexp.Regexp{
				regexp.MustCompile(`^foo:foo$`),
				regexp.MustCompile(`^bar:`),
				regexp.MustCompile(`:baz$`),
			},
			expectedOutput: "unexpected tags: qux:qux",
		},
		{
			name: "Missing tag",
			actualTags: []string{
				"foo:foo",
				"bar:bar",
				"baz:baz",
			},
			expectedTags: []*regexp.Regexp{
				regexp.MustCompile(`^foo:foo$`),
				regexp.MustCompile(`^bar:`),
				regexp.MustCompile(`:baz$`),
				regexp.MustCompile(`^qux:qux$`),
			},
			expectedOutput: "missing tags: ^qux:qux$",
		},
		{
			name: "Several mismatches",
			actualTags: []string{
				"foo:foo",
				"bar:bar",
				"baz:baz",
				"qux:qux",
				"quux:quux",
			},
			expectedTags: []*regexp.Regexp{
				regexp.MustCompile(`^foo:foo$`),
				regexp.MustCompile(`^bar:`),
				regexp.MustCompile(`:baz$`),
				regexp.MustCompile(`^corge:`),
				regexp.MustCompile(`^grault:`),
			},
			expectedOutput: "unexpected tags: qux:qux, quux:quux\nmissing tags: ^grault:, ^corge:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := assertTags(tt.actualTags, tt.expectedTags)
			if output != nil || tt.expectedOutput != "" {
				assert.EqualError(t, output, tt.expectedOutput)
			}
		})
	}
}
