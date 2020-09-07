package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateTree(root string) model.Tree {
	root, _ = filepath.Abs(root)
	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
	}

	// Create root node
	rootNode := createFileTreeNode(root, stat)
	tree := model.Tree{Root: &rootNode}

	if stat.IsDir() {
		rootNode.Children = getChildren(root)
	} else {
		rootNode.Children = make([]*model.Node, 0)
	}

	return tree
}

func getChildren(parentPath string) []*model.Node {
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

		if stat.IsDir()  {
			childNode.Children = getChildren(absPath) // parallel
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

	// TODO File type!!!
	file := model.File{Name: fileInfo.Name(), FullPath: filePath, IsDir: fileInfo.IsDir(), Size: fileInfo.Size(),
		               Created: timeStat.BirthTime(), Modified: fileInfo.ModTime(), Accessed: timeStat.AccessTime()}
	treeNode := model.Node{ Element: file }
	return treeNode

}
