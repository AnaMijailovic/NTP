"
Main window of the application, used for showing tree view and charts.
"
Class {
	#name : #ArfMainWin,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'menu',
		'pathInput',
		'chooseFileBtn',
		'chartBtns',
		'treeViewBtn',
		'pieChartBtn',
		'barChartBtn',
		'chartTypeBtns',
		'fileTypesBtn',
		'creationDateBtn',
		'accessedDateBtn',
		'modifiedDateBtn',
		'showBtn',
		'chart'
	],
	#category : #'Arf-presenters'
}

{ #category : #initialization }
ArfMainWin class >> defaultLayout [
	^ SpBoxLayout   newVertical 
			add:  #menu withConstraints: [ :constraints | constraints height: self toolbarHeight ];
			add: (SpBoxLayout newHorizontal 
						add: 'Load from: ' withConstraints: [ : constraints | constraints height: 25 ];
		   				add: #pathInput expand: true fill: true padding: 10) expand: false fill: false padding: 10;
					
			add: ( SpBoxLayout newHorizontal
					   add: ( SpBoxLayout newVertical 
							add: #treeViewBtn withConstraints: [ :constraints | constraints height: 25];
						   add: #barChartBtn withConstraints: [ :constraints | constraints height: 25];
			      		   add: #pieChartBtn withConstraints: [ :constraints | constraints height: 25];
					      add: ' '  withConstraints: [ :constraints | constraints height: 15];
						   add: #fileTypesBtn  withConstraints: [ :constraints | constraints height: 25];
			       		add: #creationDateBtn  withConstraints: [ :constraints | constraints height: 25];
			       		add: #accessedDateBtn withConstraints: [ :constraints | constraints height: 25];
							add: #modifiedDateBtn withConstraints: [ :constraints | constraints height: 25];
							add: ' '  withConstraints: [ :constraints | constraints height: 25];
							add: #showBtn expand: false fill: false padding: 25) 
					expand: false fill: false padding: 30;
					 		add: #chart expand: true fill: true padding: 1
							);
					
			yourself. 
]

{ #category : #specs }
ArfMainWin class >> defaultSpec [

	^ self defaultLayout 
]

{ #category : #'canvas drawing' }
ArfMainWin >> drawBarChart: canvas path: path chartType: chartType [

	| bar data xData yData xLabels |
	
	data := self getChartData: path chartType: chartType.
	xData := data values.
	"xData := xData collect: [ :each | each asString  ]."
	yData := data values.
	xLabels := data keys.

	bar := RSChart new
			container: canvas.
	bar addPlot: (ArfBarPlot new x: xData y: yData; xLabels: xLabels).
	
	bar addDecoration: (RSVerticalTick new 
		numberOfTicks: 5;
		asFloat).
	bar xlabel: 'File type'.
	bar ylabel: 'Size (KB)'.
	
	bar renderIn: canvas.
	bar canvas when: RSExtentChangedEvent do: [ bar canvas zoomToFit ].
]

{ #category : #'canvas drawing' }
ArfMainWin >> drawPieChart: canvas path: path chartType: chartType [

	| pie pieLabel data |
	
	data := ( self getChartData: path chartType: chartType ) associations.

	pie := RSPie new
		container: canvas;
		objects: data;
		slice: #value.
		
	pie sliceShape
		segmentSpacing: 0;
		externalRadius: 100;
		innerRadius: 50.
	pie sliceColor: NSScale category20.
	
	pieLabel := RSPieLabeled new.
	pieLabel labelShape fontSize: 4.
	pieLabel labelText: [ :obj | 
		obj key, '  (', (obj value * 100 / pie totalSum) rounded asString, ' % )'].
	
	pie build.
	pie canvas when: RSExtentChangedEvent do: [ pie canvas zoomToFit ].
	pie shapes @ pieLabel.

]

{ #category : #'canvas drawing' }
ArfMainWin >> drawTreeView: canvas path: path [
	| shapes tree createShape  |
	tree := self getTree: path.
	canvas color: '1ced1cccc'.
	createShape := [ :node | 
	| ellipse group label |
	ellipse := RSEllipse new
		model: node;
		color: (Color blue alpha: 0.7);
		size: 25;
		"draggable;"
			"contextMenu: self mainMenu;"
			when: RSMouseDoubleClick
			do: [ :evt | ArfFileDetailsPresenter new model: node element; openModalWithSpec ];
		yourself.	
	group := RSGroup new.
	label := RSLabel new
		text: node element name;
		fontSize: 3;
		yourself.
	group
		add: label;
		add: ellipse.
	RSVerticalLineLayout new alignCenter on: group.
	RSComposite new
		shapes: group;
		model: node;
		draggable;
		yourself ].
	shapes := tree asList collect: [ :node | createShape value: node ].
	canvas addAll: shapes.
	RSEdgeBuilder line
		canvas: canvas;
		objects: tree asList;
		connectToAll: #children.
	canvas edges pushBack.
	RSTreeLayout on: shapes.
	canvas @ RSCanvasController
]

{ #category : #initialization }
ArfMainWin >> getChartData: path chartType: chartType [
	| response content data |
	response := (ZnClient new)
   		url: 'http://localhost:8080/api/fileTypeData';
		queryAt: 'path' put: path;
		queryAt: 'chartType' put: chartType;
		ifFail: [ :ex | self inform: ex printString ];
   		get;
   		response.
content := response contents.

response isSuccess
   ifFalse: [ self inform: 'Something went wrong'].

content := response contents.

data := (NeoJSONReader on: content readStream)
   propertyNamesAsSymbols: true;
   next.

^ data
]

{ #category : #initialization }
ArfMainWin >> getTree: path [
	| response content tree |
	response := (ZnClient new)
   		url: 'http://localhost:8080/api/fileTree';
		queryAt: 'path' put: path;
   		get;
   		response.
content := response contents.

"Convert JSON to ArfTree object"
tree := (NeoJSONReader on: content readStream)
   mapInstVarsFor: ArfTreeNode;
   mapInstVarsFor: ArfFile;
   for: ArfTree do: [ :mapping | 
      (mapping mapInstVar: #root ) valueSchema: ArfTreeNode  ];
   for: ArfTreeNode do: [ :mapping | 
      (mapping mapInstVar: #element ) valueSchema: ArfFile .
      (mapping mapInstVar: #children ) valueSchema: #ArrayOfTreeNodes ];
   for: #ArrayOfTreeNodes
      customDo: [ :mapping | mapping listOfElementSchema: ArfTreeNode ];
   nextAs: ArfTree.

^ tree
]

{ #category : #initialization }
ArfMainWin >> initializePresenters [

	menu := self mainMenu.
	pathInput := ArfChooseFilePresenter new.

	showBtn := self newButton label: 'Show'.
	                          
	chart := self instantiate: SpRoassalPresenter.
   chart script: [ :canvas | self renderChart: canvas].
	
	treeViewBtn := self newRadioButton label: 'Tree View'.
	pieChartBtn := self newRadioButton label: 'Pie Chart'.
	barChartBtn := self newRadioButton label: 'Bar Chart'.
	treeViewBtn associatedRadioButtons: { pieChartBtn . barChartBtn }.
	chartBtns := Dictionary newFrom: { (#treeViewBtn -> treeViewBtn ) . 
	                                   (#pieChartBtn -> pieChartBtn ) .
												 (#barChartBtn -> barChartBtn ) }. 
	
	fileTypesBtn := self newRadioButton label: 'File Types'.
	creationDateBtn  := self newRadioButton label: 'Creation Date'.
	accessedDateBtn  := self newRadioButton label: 'Accessed Date'.
	modifiedDateBtn  := self newRadioButton label: 'Modified Date'.
	fileTypesBtn associatedRadioButtons: { creationDateBtn . accessedDateBtn. modifiedDateBtn }.
	chartTypeBtns := Dictionary newFrom: { (#fileType -> fileTypesBtn ) . 
	                                   (#createdDate  -> creationDateBtn ) .
												 (#accessedDate  -> accessedDateBtn ) .
												 (#modifiedDate -> modifiedDateBtn )}.
	showBtn action: [ self showChart ] 
]

{ #category : #initialization }
ArfMainWin >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'ARF - Analyze and Reorganzie Files';
		initialExtent: 900@650;
		centered
]

{ #category : #initialization }
ArfMainWin >> mainMenu [
	|aMenu|

	aMenu := self newMenuBar
	  addGroup: [ :group |
		group 
			addItem: [ :item |
				item 
					name: 'Reorganize';
					icon: (self iconNamed: #hierarchy);
				   action: [ ArfReorganizePresenter new openWithSpec ]];
			addItem: [ :item | 
				item
					name: 'Rename';
					icon: (self iconNamed: #edit) ;
					action: [ ArfRenamePresenter new openWithSpec ]];
			addItem: [ :item | 
				item
					name: 'Recover';
					icon: (self iconNamed: #glamorousRefresh) ;
					action: [ ArfRecoverPresenter new openWithSpec ]];
			addItem: [ :item |
				item
					name: 'Delete';
					icon: (self iconNamed: #glamorousTrash) ;
					action: [ArfDeletePresenter new openWithSpec]]]. 
	^ aMenu 
]

{ #category : #rendering }
ArfMainWin >> renderChart: canvas [ 
	canvas color: '6b6a66'.
]

{ #category : #'canvas drawing' }
ArfMainWin >> showChart [

	| str chartType enteredPath |

	str := ( chartBtns associations detect: [ :each | each value state ]) key.
	chartType := ( chartTypeBtns  associations detect: [ :each | each value state ]) key.
	[enteredPath := pathInput path ]
  			on: Error
  			do: [ :ex| ArfInformationDialog new message: ex messageText; openModalWithSpec. ^ self] .

	[str = 'treeViewBtn' ifTrue: [  chart script: [ :canvas | 
		self drawTreeView: canvas path: enteredPath]] .
	str = 'pieChartBtn' ifTrue: [ chart script: [ :canvas | self drawPieChart: canvas path: enteredPath chartType: chartType]].
	str = 'barChartBtn' ifTrue: [ chart script: [ :canvas | self drawBarChart: canvas path: enteredPath chartType: chartType]]]
		on: ConnectionTimedOut 
		do: [ :ex| ArfInformationDialog  new message: ex messageText ; openModalWithSpec. ^ self] .
]
