// Package git consists of all the utility functions
// that need to execute the git command.
package git

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
func TestGetGitOutput(t *testing.T) {
	testRepo := setupTestRepo()
	defer testRepo.TearDown()
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Current branch output test",
			args: args{
				args: []string{
					"symbolic-ref",
					"--short",
					"HEAD",
				},
			},
			want:    "master\n",
			wantErr: false,
		},
		{
			name: "Invalid args : error test",
			args: args{
				args: []string{
					"notAGitCommand",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGitOutput(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGitOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
