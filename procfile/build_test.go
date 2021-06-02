/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package procfile_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/procfile/procfile"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		build procfile.Build
		ctx   libcnb.BuildContext
	)

	it("does nothing without plan", func() {
		Expect(build.Build(ctx)).To(Equal(libcnb.BuildResult{}))
	})

	it("adds metadata to result, marks first process as default", func() {
		ctx.Plan = libcnb.BuildpackPlan{
			Entries: []libcnb.BuildpackPlanEntry{
				{
					Name: "procfile",
					Metadata: map[string]interface{}{
						"test-type-1": "test-command-1",
						"test-type-2": "test-command-2 argument",
					},
				},
			},
		}

		result := libcnb.NewBuildResult()
		result.Processes = append(result.Processes,
			libcnb.Process{
				Type:    "test-type-1",
				Command: "test-command-1",
				Default: true,
			},
			libcnb.Process{
				Type:    "test-type-2",
				Command: "test-command-2 argument",
			},
		)

		Expect(build.Build(ctx)).To(Equal(result))
	})

	it("adds metadata to result, marks web process as default", func() {
		ctx.Plan = libcnb.BuildpackPlan{
			Entries: []libcnb.BuildpackPlanEntry{
				{
					Name: "procfile",
					Metadata: map[string]interface{}{
						"test-type-1": "test-command-1",
						"web":         "test-command-2 argument",
					},
				},
			},
		}

		result := libcnb.NewBuildResult()
		result.Processes = append(result.Processes,
			libcnb.Process{
				Type:    "test-type-1",
				Command: "test-command-1",
			},
			libcnb.Process{
				Type:    "web",
				Command: "test-command-2 argument",
				Default: true,
			},
		)

		Expect(build.Build(ctx)).To(Equal(result))
	})

	context("tiny stack", func() {
		it.Before(func() {
			ctx.StackID = libpak.TinyStackID
		})

		it("adds metadata to result", func() {
			ctx.Plan = libcnb.BuildpackPlan{
				Entries: []libcnb.BuildpackPlanEntry{
					{
						Name: "procfile",
						Metadata: map[string]interface{}{
							"test-type-1": "test-command-1",
							"test-type-2": "test-command-2 argument",
						},
					},
				},
			}

			result := libcnb.NewBuildResult()
			result.Processes = append(result.Processes,
				libcnb.Process{
					Type:      "test-type-1",
					Command:   "test-command-1",
					Arguments: []string{},
					Direct:    true,
					Default:   true,
				},
				libcnb.Process{
					Type:      "test-type-2",
					Command:   "test-command-2",
					Arguments: []string{"argument"},
					Direct:    true,
				},
			)

			Expect(build.Build(ctx)).To(Equal(result))
		})

	})

}
