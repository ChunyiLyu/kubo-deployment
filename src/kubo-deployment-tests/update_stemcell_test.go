package kubo_deployment_tests_test

import (
	"fmt"
	"os/exec"

	. "github.com/jhvhs/gob-mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("UpdateStemcell", func() {
	It("should update the manifest with the given version", func() {
		bash.Source(pathToScript("update_stemcell"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})

		manifest := pathFromRoot("manifests/cfcr.yml")
		mockManifest := "/tmp/mock-cfcr.yml"
		cpCmd := exec.Command("cp", "-f", manifest, mockManifest)
		err := cpCmd.Run()
		Expect(err).ToNot(HaveOccurred())

		manifestFileMock := Mock("manifest_file", fmt.Sprintf("echo %s", mockManifest))
		ApplyMocks(bash, []Gob{manifestFileMock})

		exitCode, err := bash.Run("main", []string{"new-stemcell-version"})
		Expect(err).ToNot(HaveOccurred())
		Expect(exitCode).To(Equal(0))

		cmd := exec.Command("bosh-cli", "int", mockManifest, "--path=/stemcells/0/version")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("^new-stemcell-version\n$"))
	})

	It("should keep the order of the manifest the same", func() {
		bash.Source(pathToScript("update_stemcell"), nil)
		bash.Source("", func(string) ([]byte, error) {
			return repoDirectoryFunction, nil
		})

		manifest := pathFromRoot("manifests/cfcr.yml")
		mockManifest := "/tmp/mock-cfcr.yml"
		cpCmd := exec.Command("cp", "-f", manifest, mockManifest)
		err := cpCmd.Run()
		Expect(err).ToNot(HaveOccurred())

		manifestFileMock := Mock("manifest_file", fmt.Sprintf("echo %s", mockManifest))
		ApplyMocks(bash, []Gob{manifestFileMock})

		exitCode, err := bash.Run("main", []string{"new-stemcell-version"})
		Expect(err).ToNot(HaveOccurred())
		Expect(exitCode).To(Equal(0))

		// diff should only have 2 lines of change: the old version and the new version
		cmd := exec.Command("bash", "-c", fmt.Sprintf("diff -U 0 %s %s | grep -v '^@' | grep -v '^---' | grep -v '^+++' | wc -l", manifest, mockManifest))
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session).Should(gbytes.Say("^       2\n$"))
	})
})
