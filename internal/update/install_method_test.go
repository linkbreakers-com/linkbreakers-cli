package update

import "testing"

func TestInstallMethodForPath(t *testing.T) {
	cases := map[string]InstallMethod{
		"/opt/homebrew/Caskroom/linkbreakers/1.140.0/linkbreakers":                              InstallMethodHomebrew,
		"/home/linuxbrew/.linuxbrew/Caskroom/linkbreakers/1.140.0/linkbreakers":                 InstallMethodHomebrew,
		"/usr/local/Cellar/linkbreakers/1.140.0/bin/linkbreakers":                               InstallMethodHomebrew,
		"/usr/local/lib/node_modules/linkbreakers-cli/vendor/linkbreakers":                      InstallMethodNpm,
		`C:\Users\me\AppData\Roaming\npm\node_modules\linkbreakers-cli\vendor\linkbreakers.exe`: InstallMethodNpm,
		"/usr/local/bin/linkbreakers":                                                           InstallMethodDirect,
		"/home/me/.local/bin/linkbreakers":                                                      InstallMethodDirect,
	}
	for path, want := range cases {
		if got := installMethodForPath(path); got != want {
			t.Errorf("installMethodForPath(%q) = %q, want %q", path, got, want)
		}
	}
}
