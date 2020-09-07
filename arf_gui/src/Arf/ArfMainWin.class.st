Class {
	#name : #ArfMainWin,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'menu',
		'pathInput',
		'chartBtns',
		'treeViewBtn',
		'pieChartBtn',
		'barChartBtn',
		'chartTypeBtns',
		'fileTypesBtn',
		'creationDateBtn',
		'accessedDateBtn',
		'showBtn',
		'pie',
		'chart'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfMainWin class >> defaultLayout [
	^ SpBoxLayout   newVertical 
			add:  #menu withConstraints: [ :constraints | constraints height: self toolbarHeight ];
			add: ( SpBoxLayout newHorizontal
						 "beNotHomogeneous;"
					    add: ( SpBoxLayout newVertical 
							add: 'Load from: ' withConstraints: [ : constraints | constraints height: 30.
																								  constraints width: 30];
							add: #pathInput withConstraints: [ : constraints | constraints width: 50.
																								  constraints width: 30];
							add: #treeViewBtn withConstraints: [ :constraints | constraints height: 25];
						   add: #barChartBtn withConstraints: [ :constraints | constraints height: 25];
			      		   add: #pieChartBtn withConstraints: [ :constraints | constraints height: 25];
						   add: #fileTypesBtn  withConstraints: [ :constraints | constraints height: 25];
			       		add: #creationDateBtn  withConstraints: [ :constraints | constraints height: 25];
			       		add: #accessedDateBtn withConstraints: [ :constraints | constraints height: 25];
							add: #showBtn withConstraints: [:constraints | constraints height: 30]) expand: false fill: false padding: 30;
					 		add: #chart expand: true fill: true padding: 10
							);
					
			yourself. 
]

{ #category : #specs }
ArfMainWin class >> defaultSpec [

	^ self defaultLayout 
]

{ #category : #'canvas drawing' }
ArfMainWin >> drawBarChart: canvas [ 
	"TODO"
	canvas color: 'ed1cc0'.
]

{ #category : #'canvas drawing' }
ArfMainWin >> drawPieChart: canvas [ 
	"TODO"
	canvas color: '1ceaed'.
]

{ #category : #'canvas drawing' }
ArfMainWin >> drawTreeView: canvas [ 
	"TODO Delete this"

	| shapes demoFile |
	demoFile := ArfFile new 
						name: 'Ana';
						size: 22.
	canvas color: '1ced1cccc'.
	shapes := self class withAllSubclasses collect: [ :cls |
		| group label s methods|
		group := RSGroup new. 
		methods := cls methods collect: [ :met|
			RSBox new
				model: met;
				popupText: #selector;
				size: met linesOfCode;
				yourself ].
		RSGridLayout on: methods.
		s := RSComposite new
			color: (Color lightGray alpha: 0.4);
			shapes: methods;
			yourself.
		s extent: s extent + 10. 
		label := RSLabel new 
			text: cls;
			yourself.
		group add: label; add: s.
		RSVerticalLineLayout new alignCenter on: group.
		RSComposite new
			shapes: group;
			model: cls;
			draggable;
			popup;
			yourself  ].
	canvas addAll: shapes.
	RSEdgeBuilder line
		canvas: canvas;
		objects: self class withAllSubclasses; 
		connectFrom: #superclass.
	canvas edges pushBack.
	RSTreeLayout on: shapes.
	canvas @ RSCanvasController.
]

{ #category : #'canvas drawing' }
ArfMainWin >> drawTreeView: canvas path: path [

	| shapes tree createShape |
   	tree := self getTree: 'http://localhost:8080/api/fileTree?path=', path .

	canvas color: '1ced1cccc'.

	createShape := [ :node | 
		| ellipse group label|
		ellipse := RSEllipse new 
			model: node;
		   color:  (Color red alpha: 0.5);
		   size: 10;
		   draggable;
			yourself.
		group := RSGroup new. 
		label := RSLabel new 	
			text: node element name; 
			fontSize: 2;
			yourself.
		group add: label;
				add: ellipse.
	  	RSVerticalLineLayout new alignCenter on: group.
		RSComposite new
			shapes: group;
			model: node;
			draggable;
			yourself .]. 
		 
	shapes := tree asList collect: [ :node | 
		createShape value: node ]. 

	canvas addAll: shapes.
	RSEdgeBuilder line
		canvas: canvas;
		objects: tree asList; 
		connectToAll: #children.
	canvas edges pushBack.
	RSTreeLayout on: shapes.
	canvas @ RSCanvasController.

]

{ #category : #initialization }
ArfMainWin >> getDemo [
	| response content file |
	response := (ZnClient new)
   		url: 'http://localhost:8080/api/fileTree?path=C:\Users\anaam\go\src\github.com\AnaMijailovic\NTP\arf\cmd\root.go';
   		get;
   		response.
content := response contents.

file := (NeoJSONReader on: content readStream)
   mapInstVarsFor: ArfFile;
   nextAs: ArfFile.

^ self inform: file name
"^ self inform: content"
]

{ #category : #initialization }
ArfMainWin >> getTree: url [
	| response content tree |
	response := (ZnClient new)
   		url: url;
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
	pathInput := self newTextInput autoAccept: true.
	showBtn := self newButton label: 'Show'.
	
	chart := self instantiate: SpRoassalPresenter.
	pie := self instantiate: SpRoassalPresenter.
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
	fileTypesBtn associatedRadioButtons: { creationDateBtn . accessedDateBtn }.
	chartTypeBtns := Dictionary newFrom: { (#fileTypesBtn -> fileTypesBtn ) . 
	                                   (#creationDateBtn  -> creationDateBtn ) .
												 (#accessedDateBtn  -> accessedDateBtn ) }.
	showBtn action: [ self showChart ] 
]

{ #category : #initialization }
ArfMainWin >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Arf';
		initialExtent: 900@700;
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
					name: 'File';
					icon: (self iconNamed: #openIcon);
					action: [ self inform: 'Not implemented']];
		   addItem: [ :item | 
				item 
					name: 'Get Hello World';
					action: [ self getDemo ]];
			addItem: [ :item |
				item 
					name: 'Reorganize'];
			addItem: [ :item | 
				item
					name: 'Rename';
					action: [ self openRoassal ]];
			addItem: [ :item |
				item
					name: 'Delete';
					action: [ArfDeletePresenter new openWithSpec]]]. 
	^ aMenu 
]

{ #category : #initialization }
ArfMainWin >> openRoassal [
	RSShapeExamples  new example35Lines open
]

{ #category : #rendering }
ArfMainWin >> renderChart: canvas [ 
	canvas color: '6b6a66'.
]

{ #category : #'canvas drawing' }
ArfMainWin >> showChart [
	| str enteredPath |
	str := ( chartBtns   associations detect: [ :each | each value state ]) key.
	enteredPath := pathInput text.
	"self inform: ( 'Entered path: ', enteredPath )."

	str = 'treeViewBtn' ifTrue: [    chart script: [ :canvas | 
		self drawTreeView: canvas path: enteredPath]] .
	str = 'pieChartBtn' ifTrue: [ chart script: [ :canvas | self drawPieChart: canvas. ]].
	str = 'barChartBtn' ifTrue: [ chart script: [ :canvas | self drawBarChart: canvas. ]]

]
