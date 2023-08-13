package fs

import (
	"fmt"
	"github.com/heyvito/goup/models"
	"os"
	"path/filepath"
)

func DirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if stat.IsDir() {
		return true, nil
	}

	return false, fmt.Errorf("%s: exists and is not a directory", path)
}

func CurrentVersion(c models.Context) (*models.InstalledVersion, error) {
	currentPath := filepath.Join(c.InstallPath, "current")
	stat, err := os.Lstat(currentPath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	if stat.Mode()&os.ModeSymlink == 0 {
		return nil, fmt.Errorf("current path is not a symlink")
	}

	target, err := os.Readlink(currentPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading link to current installation: %w", err)
	}

	i := &models.InstalledVersion{
		Path:    target,
		RawName: filepath.Base(target),
		Status:  "",
		Ok:      false,
	}

	stat, err = os.Stat(filepath.Join(target, "bin", "go"))
	if os.IsNotExist(err) {
		i.Status = "go binary not found on current installation directory"
		return i, nil
	} else if err != nil {
		i.Status = fmt.Sprintf("failed reading go file status: %s", err)
		return i, nil
	}

	if stat.IsDir() {
		i.Status = "Binary is not a file"
		return i, nil
	}

	if stat.Mode()&0111 == 0 {
		i.Status = "Binary is not executable"
		return i, nil
	}

	i.Ok = true
	return i, nil
}

func InstalledVersions(c models.Context) ([]models.InstalledVersion, error) {
	versionsDir := filepath.Join(c.InstallPath, "versions")

	ok, _ := DirExists(versionsDir)
	if !ok {
		return nil, nil
	}

	list, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed listing versions at %s: %w", versionsDir, err)
	}

	var versions []models.InstalledVersion

	for _, node := range list {
		if !node.IsDir() {
			continue
		}

		v := models.InstalledVersion{
			Path:    filepath.Join(c.InstallPath, "versions", node.Name()),
			RawName: filepath.Base(node.Name()),
			Status:  "",
			Ok:      false,
		}

		goPath := filepath.Join(versionsDir, node.Name(), "bin", "go")
		s, err := os.Stat(goPath)
		if os.IsNotExist(err) {
			v.Status = "No binary found. Orphaned directory?"
			versions = append(versions, v)
			continue
		} else if err != nil {
			v.Status = fmt.Sprintf("Failed checking Go binary: %s", err)
			versions = append(versions, v)
			continue
		}

		if s.IsDir() {
			v.Status = "Binary is not a file"
			versions = append(versions, v)
			continue
		}

		if s.Mode()&0111 == 0 {
			v.Status = "Binary is not executable"
			versions = append(versions, v)
			continue
		}

		v.Ok = true
		versions = append(versions, v)
	}

	return versions, nil
}
