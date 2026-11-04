package builder

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDetectSourceLanguage(t *testing.T) {
	write := func(dir, name string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("# marker"), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", name, err)
		}
	}

	goOnly := t.TempDir()
	write(goOnly, "go.mod")
	if lang, err := detectSourceLanguage(goOnly); err != nil || lang != SourceLangGo {
		t.Fatalf("expected go, got %q (err=%v)", lang, err)
	}

	rustOnly := t.TempDir()
	write(rustOnly, "Cargo.toml")
	if lang, err := detectSourceLanguage(rustOnly); err != nil || lang != SourceLangRust {
		t.Fatalf("expected rust, got %q (err=%v)", lang, err)
	}

	both := t.TempDir()
	write(both, "go.mod")
	write(both, "Cargo.toml")
	if _, err := detectSourceLanguage(both); err == nil {
		t.Fatalf("expected error for ambiguous source, got nil")
	}

	neither := t.TempDir()
	if _, err := detectSourceLanguage(neither); err == nil {
		t.Fatalf("expected error for missing source, got nil")
	}
}

func TestCargoBinaryName(t *testing.T) {
	cases := []struct {
		name     string
		manifest string
		want     string
		wantErr  bool
	}{
		{
			name: "package name only",
			manifest: `[package]
name = "my-plugin"
version = "0.1.0"
`,
			want: "my-plugin",
		},
		{
			name: "explicit bin wins",
			manifest: `[package]
name = "my-plugin"

[[bin]]
name = "custom-bin"
path = "src/main.rs"
`,
			want: "custom-bin",
		},
		{
			name: "comments and spacing",
			manifest: `# top comment
[package] # trailing comment
name   =   "spaced"   # inline
version = "0.1.0"
`,
			want: "spaced",
		},
		{
			name:     "single quotes",
			manifest: "[package]\nname = 'quoted'\n",
			want:     "quoted",
		},
		{
			name:     "missing name",
			manifest: "[package]\nversion = \"0.1.0\"\n",
			wantErr:  true,
		},
		{
			name:     "empty manifest",
			manifest: "",
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			manifestPath := filepath.Join(dir, "Cargo.toml")
			if err := os.WriteFile(manifestPath, []byte(tc.manifest), 0644); err != nil {
				t.Fatalf("failed to write manifest: %v", err)
			}
			got, err := cargoBinaryName(manifestPath)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestCompileRustSource(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("cargo not installed")
	}

	dir := t.TempDir()
	manifest := `[package]
name = "nagare-test-plugin"
version = "0.1.0"
edition = "2021"
`
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write Cargo.toml: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	mainRs := "fn main() { println!(\"hello from test plugin\"); }\n"
	if err := os.WriteFile(filepath.Join(dir, "src", "main.rs"), []byte(mainRs), 0644); err != nil {
		t.Fatalf("failed to write main.rs: %v", err)
	}

	binFile := filepath.Join(dir, "pkg", "bin", "nagare-test-plugin")
	if err := compileRustSource(dir, binFile); err != nil {
		t.Fatalf("compileRustSource failed: %v", err)
	}

	info, err := os.Stat(binFile)
	if err != nil {
		t.Fatalf("expected staged binary at %s: %v", binFile, err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected non-empty staged binary")
	}

	// The staged binary must actually run.
	out, err := exec.Command(binFile).Output()
	if err != nil {
		t.Fatalf("staged binary failed to run: %v", err)
	}
	if string(out) != "hello from test plugin\n" {
		t.Fatalf("unexpected binary output: %q", out)
	}
}

func TestCompileGoSource(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not installed")
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/testplugin\n\ngo 1.27.0\n"), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}
	mainGo := "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hello from go test\") }\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	binFile := filepath.Join(dir, "pkg", "bin", "testplugin")
	if err := compileGoSource(dir, binFile); err != nil {
		t.Fatalf("compileGoSource failed: %v", err)
	}

	if _, err := os.Stat(binFile); err != nil {
		t.Fatalf("expected staged binary at %s: %v", binFile, err)
	}
}

func TestBuildRustPluginEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("cargo not installed")
	}

	dir := t.TempDir()
	metadata := `{
		"package_name": "nagare.test.rust-plugin",
		"name": "Rust Test Plugin",
		"version": "1.0.0",
		"api_version": "1.0.0",
		"author": "tester"
	}`
	if err := os.WriteFile(filepath.Join(dir, "metadata.json"), []byte(metadata), 0644); err != nil {
		t.Fatalf("failed to write metadata.json: %v", err)
	}
	manifest := "[package]\nname = \"nagare-test-rust-plugin\"\nversion = \"0.1.0\"\nedition = \"2021\"\n"
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte(manifest), 0644); err != nil {
		t.Fatalf("failed to write Cargo.toml: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src", "main.rs"), []byte("fn main() {}\n"), 0644); err != nil {
		t.Fatalf("failed to write main.rs: %v", err)
	}

	t.Chdir(dir)
	Build("")

	if _, err := os.Stat(filepath.Join(dir, "nagare.test.rust-plugin.nagare_plugin")); err != nil {
		t.Fatalf("expected packed .nagare_plugin archive: %v", err)
	}
}
