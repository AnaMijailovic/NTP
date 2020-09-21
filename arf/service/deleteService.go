package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"os"
)

func DeleteFiles(deleteData *model.DeleteData) *int {
	tree := CreateTree(deleteData.Path, deleteData.Recursive)
	filesDeleted := 0
	filesDel := &filesDeleted
	filesDel = postorderDelete(tree.Root, deleteData, &filesDeleted)
	return filesDel
}

func postorderDelete(node *model.Node, deleteData *model.DeleteData, filesDeleted *int) *int {

	for _, file := range node.Children {
		filesDeleted = postorderDelete(file, deleteData, filesDeleted)
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
		} else {
			fmt.Println("Removed: ", file.FullPath)
			*filesDeleted++
		}
	}

	return filesDeleted
}
