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

package main

import (
	"testing"

	"github.com/buildpacks/libbuildpack/v2/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/v2/build"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/cloudfoundry/procfile-cnb/procfile"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("always passes", func() {
			g.Expect(b(f.Build)).To(gomega.Equal(build.SuccessStatusCode))
		})

		it("contributes to build plan", func() {
			f.AddPlan(buildpackplan.Plan{
				Name: procfile.Dependency,
				Metadata: buildpackplan.Metadata{
					"test-type-1": "test-command-1",
					"test-type-2": "test-command-2",
				},
			})

			g.Expect(b(f.Build)).To(gomega.Equal(build.SuccessStatusCode))
			g.Expect(f.Plans).To(gomega.Equal(buildpackplan.Plans{
				Entries: []buildpackplan.Plan{
					{
						Name: procfile.Dependency,
						Metadata: buildpackplan.Metadata{
							"test-type-1": "test-command-1",
							"test-type-2": "test-command-2",
						},
					},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
