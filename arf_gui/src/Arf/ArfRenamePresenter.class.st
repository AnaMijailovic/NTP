"
Presenter for rename functionality.
"
Class {
	#name : #ArfRenamePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'pathInput',
		'recursiveCheckBox',
		'randomNamesRadioBtn',
		'removePartRadioBtn',
		'removePartInput',
		'replaceRadioBtn',
		'replaceInput',
		'replaceWithInput',
		'patternRadioBtn',
		'patternInput',
		'renameTypeBtns',
		'submitBtn'
	],
	#category : #'Arf-presenters'
}

{ #category : #initialization }
ArfRenamePresenter class >> defaultLayout [
	
	^ SpBoxLayout   newVertical 
			add: '' withConstraints: [ :c | c height: 25];
			add: 'Path: ' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: #randomNamesRadioBtn withConstraints: [ :c | c height: 45];
			add: #removePartRadioBtn withConstraints: [ :c | c height: 45];
			add: #removePartInput withConstraints: [ :c | c height: 45];
			add: #replaceRadioBtn withConstraints: [ :c | c height: 45];
			add: #replaceInput;
			add: #replaceWithInput withConstraints: [ :c | c height: 45];
			add: #patternRadioBtn withConstraints: [ :c | c height: 45];
			add: #patternInput;
			add: #submitBtn withConstraints: [ :c | c height: 40];
			yourself.
]

{ #category : #specs }
ArfRenamePresenter class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfRenamePresenter >> initializePresenters [
	pathInput := ArfChooseFilePresenter new.
	
	recursiveCheckBox := self newCheckBox label: 'Recursive'.
	
	randomNamesRadioBtn := self newRadioButton label: 'Random names'.

	removePartRadioBtn := self newRadioButton label: 'Remove part from existing name'.
	replaceRadioBtn := self newRadioButton label: 'Replace part of the existing name'.
	patternRadioBtn := self newRadioButton label: 'Use pattern:' .
	randomNamesRadioBtn associatedRadioButtons: { removePartRadioBtn . replaceRadioBtn .
	                                              patternRadioBtn }.
	renameTypeBtns  := Dictionary newFrom: { (#random -> randomNamesRadioBtn ) .  
												       (#remove-> removePartRadioBtn ) .
												       (#replace -> replaceRadioBtn ) .
												       (#pattern -> patternRadioBtn )}. 
		
	removePartInput := self newTextInput autoAccept: true .
	removePartInput placeholder: 'Enter part of the name to be removed' .
	replaceInput := self newTextInput autoAccept: true .
	replaceInput placeholder: 'Enter part of the name to be replaced' .
	replaceWithInput := self newTextInput autoAccept: true .
	replaceWithInput placeholder: 'Replace with' .	
   patternInput := self newTextInput autoAccept: true.	
	patternInput placeholder: 'Enter pattern' .

	
	submitBtn := self newButton label: 'Rename' .
	submitBtn action: [ self submitRename ] .
]

{ #category : #initialization }
ArfRenamePresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Rename';
		initialExtent: 300@500; "default size"
		centered
]

{ #category : #accessing }
ArfRenamePresenter >> pattern [
	| pattern renameBy |

	pattern := ''.
	renameBy  := ( renameTypeBtns associations detect: [ :each | each value state ]) key.
	renameBy = 'pattern' ifTrue: [ pattern  := patternInput text. 
		                            pattern = '' ifTrue: [self error: 'You must enter a pattern']].
	
	"TODO Validate pattern"
	^ pattern
]

{ #category : #accessing }
ArfRenamePresenter >> removePart [
	| removePart renameBy |

	removePart := ''.
	renameBy  := ( renameTypeBtns associations detect: [ :each | each value state ]) key.
	renameBy = 'remove' ifTrue: [ removePart  := removePartInput text].
	renameBy = 'replace' ifTrue: [ removePart  := replaceInput text].
	
	^ removePart
]

{ #category : #accessing }
ArfRenamePresenter >> replaceWith [
	| replaceWith renameBy |

	replaceWith := ''.
	renameBy  := ( renameTypeBtns associations detect: [ :each | each value state ]) key.
	renameBy = 'replace' ifTrue: [ replaceWith  := replaceWithInput text].
	
	"TODO Validate string"
	^ replaceWith
]

{ #category : #submit }
ArfRenamePresenter >> submitRename [
	
	| path recursive renameBy random remove replaceWith pattern content response|
	
	recursive := recursiveCheckBox state asString.
	
	renameBy  := ( renameTypeBtns associations detect: [ :each | each value state ]) key.
	random  := renameBy = 'random'. 
	remove := self removePart.
	replaceWith := self replaceWith.
	
	[path := pathInput path.
	 pattern := self pattern ]
  			on: Error
  			do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
	 
	[response := (ZnClient new)
   		url: 'http://localhost:8080/api/rename';
		queryAt: 'path' put: path;
		queryAt: 'recursive' put: recursive;
		queryAt: 'random' put: random;
		queryAt: 'remove' put: remove;
		queryAt: 'replaceWith' put: replaceWith ;
		queryAt: 'pattern' put: pattern;
   		put;
   		response.]
		on: ConnectionTimedOut 
		do: [ :ex| ArfInformationDialog  new message: ex messageText; openModalWithSpec. ^ self] .
			
   content := response contents. 
	ArfInformationDialog new message: content asString; openModalWithSpec.


]
