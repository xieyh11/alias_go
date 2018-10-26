package main

import "./buildtree"

func update(csvFile string, filePrefix string, nlpUsing int) {
	mapTree := buildtree.LoadTree(filePrefix)
	strMessage := buildtree.LoadMessage(filePrefix)

	names := buildtree.ReadCsv(csvFile, true, true)
	buildtree.UpdateMapTree(names, mapTree)
	buildtree.UpdateStrMessage(names, mapTree, strMessage, nlpUsing)

	buildtree.DumpTree(mapTree, filePrefix)
	buildtree.DumpMessage(strMessage, filePrefix)
}
