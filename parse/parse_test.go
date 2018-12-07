package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type TestRepo struct {
	Remote   string
	TearDown func()
}

func setupTestRepo() *TestRepo {
	pwd, _ := os.Getwd()
	oldEnv := make(map[string]string)
	overrideEnv := func(name, value string) {
		oldEnv[name] = os.Getenv(name)
		os.Setenv(name, value)
	}

	remotePath := filepath.Join(pwd, "..", "test")
	home, err := ioutil.TempDir("", "test-repo")
	if err != nil {
		panic(err)
	}

	overrideEnv("HOME", home)
	overrideEnv("XDG_CONFIG_HOME", "")
	overrideEnv("XDG_CONFIG_DIRS", "")

	targetPath := filepath.Join(home, "test")
	cmd := exec.Command("git", "clone", remotePath, targetPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("error running git clone : %s\n%s", err, output))
	}

	if err = os.Chdir(targetPath); err != nil {
		panic(err)
	}

	tearDown := func() {
		if err := os.Chdir(pwd); err != nil {
			panic(err)
		}
		for name, value := range oldEnv {
			os.Setenv(name, value)
		}
		if err = os.RemoveAll(home); err != nil {
			panic(err)
		}
	}

	return &TestRepo{Remote: remotePath, TearDown: tearDown}
}
func Test_parse(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ssh",
			args: args{
				line: "origin  git@github.com:aswinmprabhu/github-pr-cli.git (fetch)",
			},
			want: "aswinmprabhu/github-pr-cli",
		},
		{
			name: "https",
			args: args{
				line: "origin    https://github.com/aswinmprabhu/github-pr-cli (fetch)",
			},
			want: "aswinmprabhu/github-pr-cli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parse(tt.args.line); got != tt.want {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemote(t *testing.T) {
	testRepo := setupTestRepo()
	defer testRepo.TearDown()
	newHTTPSRemoteName := "httpsUpstream"
	newHTTPSRemoteURL := "https://github.com/aswinmprabhu/test"
	if err := exec.Command("git", "remote", "add", newHTTPSRemoteName, newHTTPSRemoteURL).Run(); err != nil {
		panic(fmt.Errorf("Error adding remotes : %v", err))
	}
	newSSHRemoteName := "sshUpstream"
	newSSHRemoteURL := "git@github.com:aswinmprabhu/test.git"
	if err := exec.Command("git", "remote", "add", newSSHRemoteName, newSSHRemoteURL).Run(); err != nil {
		panic(fmt.Errorf("Error adding remotes : %v", err))
	}
	type args struct {
		remoteName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "https test",
			args: args{
				remoteName: "httpsUpstream",
			},
			want:    "aswinmprabhu/test",
			wantErr: false,
		},
		{
			name: "ssh test",
			args: args{
				remoteName: "sshUpstream",
			},
			want:    "aswinmprabhu/test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Remote(tt.args.remoteName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Remote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrentBranch(t *testing.T) {
	testRepo := setupTestRepo()
	defer testRepo.TearDown()
	newBranchName := "httpsUpstream"
	if err := exec.Command("git", "checkout", "-b", newBranchName).Run(); err != nil {
		panic(fmt.Errorf("Error creating new branch : %v", err))
	}
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "New branch test",
			want:    newBranchName,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CurrentBranch()
			if (err != nil) != tt.wantErr {
				t.Errorf("CurrentBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CurrentBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}
