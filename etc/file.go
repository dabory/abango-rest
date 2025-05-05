// Author : Eric Kim
// Build Date : 23 Jul 2018  Last Update 02 Aug 2018
// End-Agent for Passcon Multi OS go binding with Windows, MacOS, iOS, and Android
// All rights are reserved.

package etc

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

func FileListToSlice(path string) ([]fs.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, LogErr("3fdgfvaqrf", FuncNameErr(), err)
	}

	// Convert DirEntry slice to FileInfo slice
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, LogErr("3fdgfvaq", FuncNameErr(), err)
		}
		files = append(files, info)
	}

	return files, nil
}

func MkdirIfNotExists(folderPath string) error {

	// Check if the folder already exists
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755) // 0755 is the permission mode
		if err != nil {
			return LogErr("3ffae3rfe", "Error creating folder:", err)
		}
	} else if err != nil {
		return LogErr("3ffae3rfe", "Unexpected Error :", err)
	} else {
		//fmt.Println("Folder already exists.")
	}

	return nil
}

// 코드가 엉망이라 전체 정비해야됨
func StrToFile(path string, str string) error {

	if file, err := os.Create(path); err != nil {
		MyErr("FDADFADSFA", err, true)
		return err
	} else {
		file.Close()
	}

	var file, err3 = os.OpenFile(path, os.O_RDWR, 0644)
	if err3 != nil {
		fmt.Println(err3)
		return err3
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err3 = file.WriteString(str)
	if err3 != nil {
		fmt.Println(err3)
		return err3
	}

	// Save file changes.
	err3 = file.Sync()
	if err3 != nil {
		return err3
	}

	return nil
}

func FileToStrWithoutErr(filename string) string {

	var str string
	if fbytes, err := os.ReadFile(filename); err == nil {
		str = string(fbytes)
	} else {
		LogStr("ETBUYIITF", "Could Not Read a file")
		return ""
	}
	return str
}

func FileToStrSkip(filename string) (string, error) {

	var str string
	if fbytes, err := os.ReadFile(filename); err == nil {
		str = string(fbytes)
	} else {
		return "", LogErr("ETBUIITF", "", err)
	}
	return str, nil
}

func FileToStr(filename string) (string, error) {

	var str string
	if fbytes, err := os.ReadFile(filename); err == nil {
		str = string(fbytes)
	} else {
		MyErr("EPOJMDOKDSF", err, true)
		return "", err
	}
	return str, nil
}

// func (t *lo.EnvConf) FileToStruct(filename string) error {

// 	if file, err := os.Open(filename); err == nil {
// 		decoder := json.NewDecoder(file)
// 		if err = decoder.Decode(t); err != nil {
// 			MyErr("LASJLDFJ", nil, true)
// 			return err
// 		}
// 	} else {
// 		MyErr("KKOIUERJ", err, true)
// 		return err
// 	}
// 	return nil
// }

// func fileCopy(src, dst string) error { // Copy시메모리 소모 없슴.
// 	sFile, err := os.Open(src)
// 	if err != nil {
// 		return MyErr("File Open", err, false)
// 	}
// 	defer sFile.Close()

// 	eFile, err := os.Create(dst)
// 	if err != nil {
// 		return MyErr("File Create", err, false)
// 	}
// 	defer eFile.Close()

// 	_, err = io.Copy(eFile, sFile) // first var shows number of bytes
// 	if err != nil {
// 		return MyErr("File Copy", err, false).Error()
// 	}

// 	err = eFile.Sync()
// 	if err != nil {
// 		return MyErr("File Open", err, false)
// 	}
// 	return nil
// }

// 정비해야지 쓸수 있슴.
func UriToFile(uri string, filepath string) error {
	// 해당내용 재정비 해야됨.  //보안부분도 강화 해야 됨.
	// Create the file=
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
