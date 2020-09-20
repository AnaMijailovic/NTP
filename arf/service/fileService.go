package service

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type WalkFunc func(path string, info os.FileInfo, err error) error
type GetLabelFunc func(filePath string) (string, error)


// *************************************************************************
//							   CREATE TREE
// *************************************************************************

func CreateTree(root string, recursive bool) model.Tree {
	root, _ = filepath.Abs(root)
	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
	}

	// Create root node
	rootNode := createFileTreeNode(root, stat)
	tree := model.Tree{Root: &rootNode}

	if stat.IsDir() {
		rootNode.Children = getChildren(root, recursive)
	} else {
		rootNode.Children = make([]*model.Node, 0)
	}

	return tree
}

func getChildren(parentPath string, recursive bool) []*model.Node {
	// ReadDir reads the directory named by dirname
	// and returns a list of directory entries sorted by filename.
	children, err := ioutil.ReadDir(parentPath)

	if err != nil {
		fmt.Println(err)
	}

	childrenNodes := make([]*model.Node, 0)

	for _, file := range children {

		absPath, _ := filepath.Abs(parentPath + string(os.PathSeparator) + file.Name() )
		stat, err := os.Stat(absPath)

		if err != nil {
			fmt.Println(err)
		}

		childNode := createFileTreeNode(absPath, stat)
		childrenNodes = append(childrenNodes, &childNode)

		if stat.IsDir()  && recursive {
			childNode.Children = getChildren(absPath, recursive) // parallel
		} else {
			childNode.Children = make([]*model.Node, 0)
		}
	}

	return childrenNodes
}

func createFileTreeNode(filePath string, fileInfo os.FileInfo) model.Node {
	// for created and accessed dates
	// they are not provided in FileInfo
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	fileType, _ := getFileContentType(filePath)
	file := model.File{Name: fileInfo.Name(), FullPath: filePath, IsDir: fileInfo.IsDir(), Size: fileInfo.Size(),
		Created: timeStat.BirthTime(), Modified: fileInfo.ModTime(), Accessed: timeStat.AccessTime(),
		FileType: fileType }
	treeNode := model.Node{ Element: file }
	return treeNode

}

func getFileContentType(filePath string) (string, error) {
	out, _ := os.Open(filePath)
	defer out.Close()
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	// fmt.Println(contentType)
	return contentType, nil
}

// *************************************************************************
//							   CHART DATA
// *************************************************************************

func GetFileChartData(root string, chartType string) map[string]int64 {

	getLabelFunc := getLabelFunc(chartType)
	// closure
	fileTypesDict := make(map[string]int64)

	Walk(root, true, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		fileType, _ := getLabelFunc(path)
		size, ok := fileTypesDict[fileType]
		if ok {
			fileTypesDict[fileType] = size + info.Size()
		} else {
			fileTypesDict[fileType] = info.Size()
		}

		return nil
	})

	return fileTypesDict
}

func getLabelFunc(chartType string) GetLabelFunc {

	if chartType == "fileType"{
		return getFileContentType
	} else if chartType == "createdDate" {
		return getCreatedDateLabel
	} else if chartType == "modifiedDate" {
		return getModifiedDateLabel
	}else { // accessed date
		return getAccessedDateLabel
	}
}

func getCreatedDateLabel(filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	createdYear := timeStat.BirthTime().Year()

	return strconv.Itoa(createdYear), nil
}

func getModifiedDateLabel(filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	modifiedYear := timeStat.ModTime().Year()

	return strconv.Itoa(modifiedYear), nil
}

func getAccessedDateLabel(filePath string) (string, error) {
	timeStat, err := times.Stat(filePath)

	if err != nil {
		fmt.Println(err)
	}

	accessedYear := timeStat.ModTime().Year()

	return strconv.Itoa(accessedYear), nil
}

func Walk(root string, recursive bool, walkFn WalkFunc) error {

	stat, err := os.Stat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}

	if !stat.IsDir() {
		return walkFn(root, stat, nil)
	}

	// if there is no error
	// first call walkFn for the root
	walkFn(root, stat, nil)

	// ReadDir reads the directory named by dirname
	// and returns a list of directory entries sorted by filename.
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {

		abs, _ := filepath.Abs(root + string(os.PathSeparator) + file.Name() )
		stat, err := os.Stat(abs)

		if err != nil {
			return err
		}

		if stat.IsDir() && recursive {
			Walk(abs, recursive, walkFn)
		} else { // If it is a file
			walkFn(abs, stat, nil)
		}
	}

	return nil
}


// *************************************************************************
//								DELETE
// *************************************************************************

func DeleteFiles(deleteData *model.DeleteData) {
	tree := CreateTree(deleteData.Path, deleteData.Recursive)

	postorderDelete(tree.Root, deleteData)

}

func postorderDelete(node *model.Node, deleteData *model.DeleteData ) {

	for _, file := range node.Children {
		postorderDelete(file, deleteData)
	}

	file := node.Element.(model.File)
	stat, _ := os.Stat(file.FullPath)
	size := stat.Size()

	if ( deleteData.Empty && size == 0 ) ||
		(deleteData.CreatedBefore.After(file.Created) && (!file.IsDir || file.Size == 0)) ||
		(deleteData.NotAccessedAfter.After(file.Accessed) && (!file.IsDir || file.Size == 0)){
		err := os.Remove(file.FullPath)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Removed: ", file.FullPath)
	}
}

// *************************************************************************
//								REORGANIZE
// *************************************************************************

func ReorganizeFiles(reorganizeData *model.ReorganizeData) {
	recoveryFilePath := reorganizeData.Dest + string(os.PathSeparator) + "arfRecover.txt"
	// TODO Something better here?
	os.Remove(recoveryFilePath)

	Walk(reorganizeData.Src, reorganizeData.Recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFolderName := generateFolderName(path, reorganizeData)
		newFolderName = strings.Replace(newFolderName, "/", "_", -1)

		newFolderPath := reorganizeData.Dest + string(os.PathSeparator) + newFolderName
		fmt.Println("New folder path: ", newFolderPath)

		// check if directory already exists
		if _, err := os.Open(newFolderPath); err != nil {
			err := os.Mkdir(newFolderPath, 0755)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Move file
		newFilePath :=  newFolderPath + string(os.PathSeparator) + info.Name()
		err = os.Rename(path,newFilePath)

		writeRecoveryData(recoveryFilePath, path, newFilePath)

		return nil
	})

}

func writeRecoveryData(recoveryFilePath string, src string, dest string) {

	// If the file doesn't exist, create it, or append to the file
	recoveryFile, err := os.OpenFile(recoveryFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer recoveryFile.Close()

	recoveryFile.WriteString(src + "," + dest + "\n")
}

func generateFolderName(path string, reorganizeData *model.ReorganizeData) string {
	if reorganizeData.FileType {
		return generateFolderNameByType(path)
	} else if reorganizeData.FileSize != 0 {
		return generateFolderNameBySize(path, reorganizeData.FileSize)
	} else {
		return generateFolderNameByDate(path, reorganizeData.CreatedDate)
	}
}

func generateFolderNameByType(path string) string {
	fileType, _ := getFileContentType(path)
	return fileType
}

func generateFolderNameBySize(path string, step int64) string {
	var label string
	stat, _ := os.Stat(path)
	size := stat.Size() / 1000000 // MB

	groupNum := size / step

	if groupNum == 0 {
		label = strconv.FormatInt(step*groupNum, 10) + "-" + strconv.FormatInt(step*groupNum + step, 10)
	} else {
		label = strconv.FormatInt(step*groupNum + 1, 10) + "-" + strconv.FormatInt(step*groupNum + step, 10)
	}

	return label
}

func generateFolderNameByDate(path string, dateType string) string {

	timeStat, _ := times.Stat(path)
	var folderName string

	if dateType == "d" {
		folderName = timeStat.BirthTime().Format("02-01-2006")
	} else if dateType == "m" {
		folderName = timeStat.BirthTime().Format("01-2006")
	} else { // year
		folderName = timeStat.BirthTime().Format("2006")
	}

	return folderName
}

// *************************************************************************
//								RENAME
// *************************************************************************

func Rename(renameData *model.RenameData) {
	recoveryFilePath := renameData.Path + string(os.PathSeparator) + "arfRecover.txt"
	// TODO Something better here?
	os.Remove(recoveryFilePath)

	Walk(renameData.Path, renameData.Recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFileName := generateNewFileName(path, renameData)

		newFilePath := filepath.Dir(path) + string(os.PathSeparator) + newFileName

		fmt.Println("New file path: ", newFilePath)

		// Rename a file
		err = os.Rename(path, newFilePath)


		if err != nil {
			log.Fatal(err)
		}

		writeRecoveryData(recoveryFilePath, path, newFilePath)

		return nil
	})
}

func generateNewFileName(oldFilePath string, renameData *model.RenameData) string {
	oldFileName := filepath.Base(oldFilePath)
	if renameData.Random {
		randomStr, _ := GenerateRandomString(12)
		extension := filepath.Ext(oldFilePath)
		randomStr = randomStr + extension
		return  randomStr
	} else if renameData.Remove != "" {
		return strings.ReplaceAll(oldFileName, renameData.Remove, renameData.ReplaceWith)
	} else if renameData.Pattern != "" {
		newName, _ := parsePatternString(renameData.Pattern, oldFileName )
		fmt.Println("New name: ", newName)
		return newName
	}

	//TODO other cases
	return "TODO..."
}

func parsePatternString(patternStr string, oldFileName string) (string, error){
	// var newName string
	// Remove extension from file name
	ext := filepath.Ext(oldFileName)
	oldFileName = strings.Replace(oldFileName, ext, "", 1)

	for strings.Index(patternStr, "{") != -1 {
		startIndex := strings.Index(patternStr, "{")
		endIndex := strings.Index(patternStr, "}")
		tag := patternStr[startIndex+1 : endIndex]

		fmt.Println("Between: ", tag)

		if tag == "name" {
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], oldFileName, 1)

		}else if tag == "random" {
			random, _ := GenerateRandomString(12)
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], random, 1 )
		}

	}
	return patternStr + ext, nil

}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// *************************************************************************
//								RECOVER
// *************************************************************************

func Recover(recoveryFilePath string) []error {
	errs := make([]error, 0)

	file, err := os.Open(recoveryFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Split next line
		paths := strings.Split(scanner.Text(), ",")
		src, dest := paths[0], paths[1]

		// TODO Return all that failed to move
		err := moveFile(dest, src)
		if err != nil {
			errs = append(errs, err)
		}

		// Delete destination directory if it is empty
		file,_ := os.Stat(filepath.Dir(dest))
		if file.Size() == 0 {
			os.Remove(filepath.Dir(dest))
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Delete recoveryFile
	file.Close()
	err = os.Remove(recoveryFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return errs
}

func moveFile(src string, dest string) error {

	dirPath := filepath.Dir(dest)

	if _, err := os.Open(dirPath); err != nil {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	// Move file
	err := os.Rename(src, dest)

	if err != nil {
		return err
	}

	return nil

}