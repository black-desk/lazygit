package commands_test

import (
	"github.com/go-errors/errors"
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StatusMgr", func() {
	var (
		commander *FakeICommander
		gitconfig *FakeIGitConfigMgr
		StatusMgr *StatusMgr
		mgrCtx    *MgrCtx
	)

	BeforeEach(func() {
		commander = NewFakeCommander()
		gitconfig = &FakeIGitConfigMgr{}

		mgrCtx = NewFakeMgrCtx(commander, gitconfig, nil)

		StatusMgr = NewStatusMgr(mgrCtx)
	})

	Describe("CurrentBranchName", func() {
		Context("On master branch", func() {
			It("returns 'master'", func() {
				ExpectRunWithOutputCalls(commander, []ExpectedRunCall{
					{"git symbolic-ref --short HEAD", "master\n", nil},
				})

				name, displayname, err := StatusMgr.CurrentBranchName()
				Expect(name).To(Equal("master"))
				Expect(displayname).To(Equal("master"))
				Expect(err).To(BeNil())
			})
		})

		Context("symbolic-ref fails", func() {
			Context("when git branch command says we're on master", func() {
				It("falls back to 'git branch --contains'", func() {
					ExpectRunWithOutputCalls(commander, []ExpectedRunCall{
						{"git symbolic-ref --short HEAD", "", errors.New("my error")},
						{"git branch --contains", "* master\n  otherbranch\n", nil},
					})

					name, displayname, err := StatusMgr.CurrentBranchName()
					Expect(name).To(Equal("master"))
					Expect(displayname).To(Equal("master"))
					Expect(err).To(BeNil())
				})
			})

			Context("when git branch command says we're on a detached head", func() {
				It("falls back to 'git branch --contains'", func() {
					ExpectRunWithOutputCalls(commander, []ExpectedRunCall{
						{"git symbolic-ref --short HEAD", "", errors.New("my error")},
						{"git branch --contains", "* (HEAD detached at 264fc6f5)\n  otherbranch\n", nil},
					})

					name, displayname, err := StatusMgr.CurrentBranchName()
					Expect(name).To(Equal("264fc6f5"))
					Expect(displayname).To(Equal("(HEAD detached at 264fc6f5)"))
					Expect(err).To(BeNil())
				})
			})

			Context("when both commands return an error", func() {
				It("bubbles up error", func() {
					ExpectRunWithOutputCalls(commander, []ExpectedRunCall{
						{"git symbolic-ref --short HEAD", "", errors.New("my error")},
						{"git branch --contains", "", errors.New("my other error")},
					})

					name, displayname, err := StatusMgr.CurrentBranchName()
					Expect(name).To(Equal(""))
					Expect(displayname).To(Equal(""))
					Expect(err).To(MatchError("my other error"))
				})
			})
		})
	})
})
