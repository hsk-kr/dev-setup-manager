package tools

type InstalledSoftwareList struct {
	Homebrew bool
}

/*
first value contains if software is installed or not,
to update this first value, you should call the function given in the return value
*/
func CreateInstalledSoftwareList() (*InstalledSoftwareList, func()) {
	installedSoftwareList := InstalledSoftwareList{}

	updateInstalledSoftware := func() {
		installedSoftwareList.Homebrew = IsHomebrewInstalled()
	}

	updateInstalledSoftware()

	return &installedSoftwareList, updateInstalledSoftware
}
