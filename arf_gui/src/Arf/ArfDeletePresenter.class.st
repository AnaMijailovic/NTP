Class {
	#name : #ArfDeletePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'pathInput',
		'recursiveCheckBox',
		'emptyCheckBox',
		'creationCheckBox',
		'creationDateInput',
		'accessedCheckBox',
		'accessedDateInput',
		'submitBtn'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfDeletePresenter class >> defaultLayout [
	^ SpBoxLayout   newVertical 
			add: '' withConstraints: [ :c | c height: 25];
			add: 'Path: ' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: #emptyCheckBox withConstraints: [ :c | c height: 45];
			add: #creationCheckBox withConstraints: [ :c | c height: 45];
			add: #creationDateInput withConstraints: [ :c | c height: 45];
			add: #accessedCheckBox withConstraints: [ :c | c height: 45];
			add: #accessedDateInput;
			add: #submitBtn withConstraints: [ :c | c height: 40];
		   yourself.
]

{ #category : #specs }
ArfDeletePresenter class >> defaultSpec [

	^ self defaultLayout 
]

{ #category : #accessing }
ArfDeletePresenter >> createdBefore [
	| createdBefore |
	createdBefore := '01-01-0001' .
	creationCheckBox state  ifTrue: [ createdBefore := creationDateInput date printFormat: #(1 2 3 $- 2 1)] .
	^ createdBefore 
]

{ #category : #initialization }
ArfDeletePresenter >> initializePresenters [
	pathInput := self newTextInput autoAccept: true.
	
	recursiveCheckBox := self newCheckBox label: 'Recursive'.
	
	emptyCheckBox := self newCheckBox label: 'Delete empty files and folders'.
	creationCheckBox  := self newCheckBox label: 'Delete files and folders created before:'.
	accessedCheckBox := self newCheckBox label: 'Delete files and folders not accessed after:'.
												
	creationDateInput := self instantiate: SpDatePresenter.
	creationDateInput display: [ :item | item yyyymmdd ].
	
	accessedDateInput := self instantiate: SpDatePresenter.
	accessedDateInput display: [ :item | item yyyymmdd ].
	
	submitBtn := self newButton label: 'Delete'.
	submitBtn action: [ self submitDelete ].
]

{ #category : #initialization }
ArfDeletePresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Delete';
		initialExtent: 300@415;
		centered
]

{ #category : #accessing }
ArfDeletePresenter >> notAccessedAfter [
	| notAccessedAfter |
	notAccessedAfter := '01-01-0001' .
	creationCheckBox state  ifTrue: [ notAccessedAfter := creationDateInput date printFormat: #(1 2 3 $- 2 1)] .
	^ notAccessedAfter 
]

{ #category : #accessing }
ArfDeletePresenter >> path [ 
	|path isValidPath |

	path := pathInput text.
	path = '' ifTrue: [ self error: 'You must enter a path' ]. 
	
	isValidPath := path asFileReference.
	isValidPath exists ifFalse: [ self error: 'Path is invalid' ] .
	
	^ path
]

{ #category : #initialization }
ArfDeletePresenter >> submitDelete [
	| enteredPath recursive empty createdBefore notAccessedAfter response content |

	recursive := recursiveCheckBox state asString.
	empty := emptyCheckBox  state asString.
	
	createdBefore := self createdBefore  .
	notAccessedAfter := self notAccessedAfter .
	
	[enteredPath := self path ]
  			on: Error
  			do: [ :ex| ArfInformationDialog new message: ex messageText; openModalWithSpec. ^ self] .
	
	"Validate that at least one criteria is provided"
	createdBefore = '01-01-0001' ifTrue: [ 
		notAccessedAfter = '01-01-0001' ifTrue: [ 
			emptyCheckBox state ifFalse: [ 
				ArfInformationDialog new message: 'Critera not provided'; openModalWithSpec. ^ self.
				 ]
			 ]
		] .
	
	[response := (ZnClient new)
   		url: 'http://localhost:8080/api/deleteFiles';
		queryAt: 'path' put: enteredPath ;
		queryAt: 'recursive' put: recursive;
		queryAt: 'empty' put: empty;
		queryAt: 'createdBefore' put: self createdBefore ;
		queryAt: 'notAccessedAfter' put: self notAccessedAfter;
   		delete;
   		response]
		on: ConnectionTimedOut 
		do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
			
   content := response contents. 
	ArfInformationDialog new message: content asString; openModalWithSpec.
]
