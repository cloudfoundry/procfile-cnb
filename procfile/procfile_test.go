/*
 * Copyright 2018-2019 the original author or authors.
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
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/cloudfoundry/procfile-cnb/procfile"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestProcfile(t *testing.T) {
	spec.Run(t, "Procfile", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		when("NewProcfile", func() {

			it("returns false when no procfile", func() {
				_, ok, err := procfile.NewProcfile(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeFalse())
			})

			it("returns true when procfile exists", func() {
				f.AddPlan(buildpackplan.Plan{Name: procfile.Dependency})

				_, ok, err := procfile.NewProcfile(f.Build)

				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeTrue())
			})
		})

		when("ParseProcfile", func() {

			var f *test.DetectFactory

			it.Before(func() {
				f = test.NewDetectFactory(t)
			})

			it("returns false when no Procfile", func() {
				_, ok, err := procfile.ParseProcfile(f.Detect.Application, f.Detect.Logger)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeFalse())
			})

			it("returns true when Procfile exists", func() {
				test.WriteFile(t, filepath.Join(f.Detect.Application.Root, "Procfile"), "test-type: test-command")

				p, ok, err := procfile.ParseProcfile(f.Detect.Application, f.Detect.Logger)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeTrue())
				g.Expect(p).To(gomega.Equal(map[string]string{"test-type": "test-command"}))
			})
		})

		it("contributes command", func() {
			f.AddPlan(buildpackplan.Plan{
				Name: procfile.Dependency,
				Metadata: buildpackplan.Metadata{
					"test-type-1": "test-command-1",
					"test-type-2": "test-command-2",
				},
			})

			p, _, err := procfile.NewProcfile(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(p.Contribute()).To(gomega.Succeed())

			g.Expect(f.Build.Layers).To(test.HaveApplicationMetadata(layers.Metadata{
				Processes: []layers.Process{
					{Type: "test-type-1", Command: "test-command-1"},
					{Type: "test-type-2", Command: "test-command-2"},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
