package common

import (
	"os"
	"path"
	fp "path/filepath"
)

// AdjustPathSlash standardize path
// separators according to operating system
func AdjustPathSlash(path string) string {

	return fp.FromSlash(path)
}

// ChangeDir - Chdir changes the current working directory to the named directory. If there is an error, it will be of type *PathError.
func ChangeDir(dirPath string) error {

	err := os.Chdir(dirPath)

	if err != nil {
		return err
	}

	return nil
}

// CreateFile - Wrapper function for os.Create
func CreateFile(fileName string) (*os.File, error) {
	return os.Create(fileName)
}

// DeleteDirFile - Wrapper function for Remove.
// Remove removes the named file or directory.
// If there is an error, it will be of type *PathError.
func DeleteDirFile(pathFile string) error {
	return os.Remove(pathFile)
}

// DeleteDirPathAll - Wrapper function for RemoveAll
// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first
// error it encounters. If the path does not exist,
// RemoveAll returns nil (no error).
func DeleteDirPathAll(path string) error {
	return os.RemoveAll(path)
}

// DoesFileExist - Returns a boolean value
// designating whether the passed file name
// exists.
func DoesFileExist(pathFileName string) bool {

	status, _, _ := DoesFileInfoExist(pathFileName)

	return status
}

// DoesFileInfoExist - returns a boolean value indicating
// whether the path and file name passed to the function
// actually exists. Note: If the file actually exists,
// the function will return the associated FileInfo structure.
func DoesFileInfoExist(pathFileName string) (bool, os.FileInfo, error) {
	var fInfo os.FileInfo
	var err error

	if fInfo, err = os.Stat(pathFileName); os.IsNotExist(err) {
		return false, fInfo, err
	}

	return true, fInfo, nil

}

// GetAbsPathFromFilePath - Supply a string containing both
// the path file name and extension and return the path
// element.
func GetAbsPathFromFilePath(filePath string) (string, error) {

	return MakeAbsolutePath(path.Dir(filePath))

}

// GetAbsCurrDir - returns
// the absolute path of the
// current working directory
func GetAbsCurrDir() (string, error) {

	dir, err := GetCurrentDir()

	if err != nil {
		return dir, err
	}

	return MakeAbsolutePath(dir)
}

// GetCurrentDir - Wrapper function for
// Getwd(). Getwd returns a rooted path name
// corresponding to the current directory.
// If the current directory can be reached via
// multiple paths (due to symbolic links),
// Getwd may return any one of them.
func GetCurrentDir() (string, error) {
	return os.Getwd()
}

// GetExecutablePathFileName - Gets the file name
// and path of the executable that started the
// current process
func GetExecutablePathFileName() (string, error) {
	ex, err := os.Executable()

	return ex, err

}

// JoinPathsAdjustSeparators - Joins two
// path strings and standardizes the
// path separators according to the
// current operating system.
func JoinPathsAdjustSeparators(p1 string, p2 string) string {

	return fp.FromSlash(JoinPaths(p1, p2))

}

// JoinPaths - correctly joins 2-paths
func JoinPaths(p1 string, p2 string) string {

	return path.Join(p1, p2)

}

// MakeAbsolutePath - Supply a relative path or any path
// string and resolve that path to an Absolute Path.
func MakeAbsolutePath(relPath string) (string, error) {

	path, err := fp.Abs(relPath)

	if err != nil {
		return "Invalid Path!", err
	}

	return path, err
}

// MakeDirAll - creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error. The permission bits perm
// are used for all directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func MakeDirAll(dirPath string) error {
	var ModePerm os.FileMode = 0777
	return os.MkdirAll(dirPath, ModePerm)
}

// MakeDir - Makes a directory. Returns
// boolean value of false plus error if
// the operation fails. If successful,
// the function returns true.
func MakeDir(dirPath string) (bool, error) {
	var ModePerm os.FileMode = 0777
	err := os.Mkdir(dirPath, ModePerm)

	if err != nil {
		return false, err
	}

	return true, nil
}
