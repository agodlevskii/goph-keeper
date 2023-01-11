package configs

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCfg struct {
	Name string `json:"name" yaml:"name" env:"TEST_NAME" envDefault:"test"`
	Age  int    `json:"age" yaml:"age" env:"TEST_AGE"`
}

func TestUpdateConfigFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *testCfg
		want    *testCfg
		wantErr bool
	}{
		{
			name: "Default env config",
			cfg:  &testCfg{},
			want: &testCfg{Name: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateConfigFromEnv(tt.cfg)
			assert.Equal(t, tt.want, tt.cfg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUpdateConfigFromFile(t *testing.T) {
	type args struct {
		cfg   *testCfg
		fCfg  *testCfg
		fPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *testCfg
		wantErr bool
	}{
		{
			name:    "Config filepath is missing",
			wantErr: true,
		},
		{
			name: "Default config",
			args: args{
				cfg:   &testCfg{},
				fCfg:  &testCfg{},
				fPath: "test_cfg.json",
			},
			want: &testCfg{},
		},
		{
			name: "Wrong config extension",
			args: args{
				cfg:   &testCfg{},
				fCfg:  &testCfg{Name: "test1", Age: 1},
				fPath: "test_cfg.txt",
			},
			want:    &testCfg{Name: "", Age: 0},
			wantErr: true,
		},
		{
			name: "Correct JSON config",
			args: args{
				cfg:   &testCfg{},
				fCfg:  &testCfg{Name: "test1", Age: 1},
				fPath: "test_cfg.json",
			},
			want: &testCfg{Name: "test1", Age: 1},
		},
		{
			name: "Existing JSON config",
			args: args{
				cfg:   &testCfg{Age: 2},
				fCfg:  &testCfg{Name: "test1", Age: 1},
				fPath: "test_cfg.json",
			},
			want: &testCfg{Name: "test1", Age: 2},
		},
		{
			name: "Existing yaml config",
			args: args{
				cfg:   &testCfg{Age: 2},
				fCfg:  &testCfg{Name: "test1", Age: 3},
				fPath: "test_cfg.yml",
			},
			want: &testCfg{Name: "test1", Age: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fPath := generateConfigFileName(t, tt.args.fPath)
			var f *os.File
			if fPath != "" {
				f = setupFileConfig(t, fPath, tt.args.fCfg)
			}

			err := UpdateConfigFromFile(tt.args.cfg, tt.args.fCfg, fPath)
			assert.Equal(t, tt.want, tt.args.cfg)
			assert.Equal(t, tt.wantErr, err != nil)

			if fPath != "" {
				if err = f.Close(); err != nil {
					t.Fatal(err)
				}
				cleanFileConfig(t, fPath)
			}
		})
	}
}

func setupFileConfig(t *testing.T, filename string, cfg *testCfg) *os.File {
	path, err := getConfigsDirPath()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatal(err)
	}
	return f
}

func cleanFileConfig(t *testing.T, filename string) {
	path, err := getConfigsDirPath()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(filepath.Join(path, filename)); err != nil {
		t.Fatal(err)
	}
}

func generateConfigFileName(t *testing.T, fName string) string {
	if fName == "" {
		return ""
	}

	r, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		t.Fatal(err)
	}
	return strconv.Itoa(int(r.Int64())) + fName
}
