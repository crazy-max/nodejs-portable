package fs

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// File is an open file on a file system.
type File interface {
	io.Reader
	io.Writer
	io.Closer

	Fd() uintptr
	Readdirnames(n int) ([]string, error)
	Readdir(int) ([]os.FileInfo, error)
	Seek(int64, int) (int64, error)
	Stat() (os.FileInfo, error)
}

func fixpath(name string) string {
	abspath, err := filepath.Abs(name)
	if err == nil {
		return UncPath(abspath)
	}
	return name
}

// UncPath converts an absolute Windows path to a UNC long path.
func UncPath(s string) string {
	isAbsWinDrive := regexp.MustCompile(`^[a-zA-Z]\:\\`)

	// UNC can NOT use "/", so convert all to "\"
	s = strings.Replace(s, `/`, `\`, -1)

	// If prefix is "\\", we already have a UNC path or server.
	if strings.HasPrefix(s, `\\`) {
		// If already long path, just keep it
		if strings.HasPrefix(s, `\\?\`) {
			return s
		}
		// Trim "\\" from path and add UNC prefix.
		return `\\?\UNC\` + strings.TrimPrefix(s, `\\`)
	}

	if isAbsWinDrive.MatchString(s) {
		return `\\?\` + s
	}

	return s
}

func RemoveUnc(s string) string {
	return strings.TrimLeft(s, `\\?\`)
}

// Open opens a file for reading.
func Open(name string) (File, error) {
	return os.OpenFile(fixpath(name), os.O_RDONLY, 0)
}

// ClearCache syncs and then removes the file's content from the OS cache.
func ClearCache(f File) error {
	return nil
}

// Chmod changes the mode of the named file to mode.
func Chmod(name string, mode os.FileMode) error {
	return os.Chmod(fixpath(name), mode)
}

// Mkdir creates a new directory with the specified name and permission bits.
// If there is an error, it will be of type *PathError.
func Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(fixpath(name), perm)
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(fixpath(path), perm)
}

// Readlink returns the destination of the named symbolic link.
// If there is an error, it will be of type *PathError.
func Readlink(name string) (string, error) {
	return os.Readlink(fixpath(name))
}

// Remove removes the named file or directory.
// If there is an error, it will be of type *PathError.
func Remove(name string) error {
	return os.Remove(fixpath(name))
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters.  If the path does not exist, RemoveAll
// returns nil (no error).
func RemoveAll(path string) error {
	return os.RemoveAll(fixpath(path))
}

// Rename renames (moves) oldpath to newpath.
// If newpath already exists, Rename replaces it.
// OS-specific restrictions may apply when oldpath and newpath are in different directories.
// If there is an error, it will be of type *LinkError.
func Rename(oldpath, newpath string) error {
	return os.Rename(fixpath(oldpath), fixpath(newpath))
}

// Symlink creates newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.
func Symlink(oldname, newname string) error {
	return os.Symlink(fixpath(oldname), fixpath(newname))
}

// Stat returns a FileInfo structure describing the named file.
// If there is an error, it will be of type *PathError.
func Stat(name string) (os.FileInfo, error) {
	return os.Stat(fixpath(name))
}

// Lstat returns the FileInfo structure describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link.  Lstat makes no attempt to follow the link.
// If there is an error, it will be of type *PathError.
func Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(fixpath(name))
}

// Create creates the named file with mode 0666 (before umask), truncating
// it if it already exists. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode
// O_RDWR.
// If there is an error, it will be of type *PathError.
func Create(name string) (*os.File, error) {
	return os.Create(fixpath(name))
}

// OpenFile is the generalized open call; most users will use Open
// or Create instead.  It opens the named file with specified flag
// (O_RDONLY etc.) and perm, (0666 etc.) if applicable.  If successful,
// methods on the returned File can be used for I/O.
// If there is an error, it will be of type *PathError.
func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(fixpath(name), flag, perm)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order, which makes the output deterministic but means that for very
// large directories Walk can be inefficient.
// Walk does not follow symbolic links.
func Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(fixpath(root), walkFn)
}

// Join joins any number of path elements into a single path, adding a
// separating slash if necessary. The result is Cleaned; in particular,
// all empty strings are ignored.
func Join(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return fixpath(path.Clean(strings.Join(elem[i:], "/")))
		}
	}
	return ""
}

// CopyFile copies file source to destination dest.
func CopyFile(src, dst string) (err error) {
	in, err := Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := Stat(src)
	if err != nil {
		return
	}

	err = Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	err = MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := Join(src, entry.Name())
		dstPath := Join(dst, entry.Name())
		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// CreateSubfolder creates a folder recusively
func CreateSubfolder(path string) error {
	if _, err := Stat(path); os.IsNotExist(err) {
		err = MkdirAll(path, 777)
		if err != nil {
			return err
		}
	}
	return nil
}

// FormatWinPath replaces / with \
func FormatWinPath(path string) string {
	return strings.Replace(path, "/", "\\", -1)
}

// Unzip a file in a destination path
// https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		thePath := Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			MkdirAll(thePath, f.Mode())
		} else {
			MkdirAll(filepath.Dir(thePath), f.Mode())
			f, err := OpenFile(thePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
