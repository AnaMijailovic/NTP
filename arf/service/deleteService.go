package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"os"
)

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
