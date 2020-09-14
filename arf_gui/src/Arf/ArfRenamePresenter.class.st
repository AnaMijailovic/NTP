Class {
	#name : #ArfRenamePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'pathInput',
		'recursiveCheckBox',
		'randomNamesRadioBtn',
		'addPrefixRadioBtn',
		'prefixInput',
		'addSufixRadioBtn',
		'sufixInput',
		'removePartRadioBtn',
		'removePartInput',
		'replaceRadioBtn',
		'replaceInput',
		'replaceWithInput',
		'renameTypeBtns',
		'submitBtn'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfRenamePresenter class >> defaultLayout [
	
	^ SpBoxLayout   newVertical 
			add: 'Path: ' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: #randomNamesRadioBtn withConstraints: [ :c | c height: 45];
			add: #addPrefixRadioBtn withConstraints: [ :c | c height: 45];
			add: #prefixInput withConstraints: [ :c | c height: 45];
			add: #addSufixRadioBtn withConstraints: [ :c | c height: 45];
			add: #sufixInput withConstraints: [ :c | c height: 45];
			add: #removePartRadioBtn withConstraints: [ :c | c height: 45];
			add: #removePartInput withConstraints: [ :c | c height: 45];
			add: #replaceRadioBtn withConstraints: [ :c | c height: 45];
			add: #replaceInput;
			add: #replaceWithInput;
			add: #submitBtn withConstraints: [ :c | c height: 35];
			yourself.
]

{ #category : #specs }
ArfRenamePresenter class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfRenamePresenter >> initializePresenters [
	pathInput := self newTextInput autoAccept: true.
	
	recursiveCheckBox := self newCheckBox label: 'Recursive'.
	
	randomNamesRadioBtn := self newRadioButton label: 'Random names'.
	addPrefixRadioBtn := self newRadioButton label: 'Add prefix to existing name'.
	addSufixRadioBtn := self newRadioButton label: 'Add sufix to existing name'.
	removePartRadioBtn := self newRadioButton label: 'Remove part from existing name'.
	replaceRadioBtn := self newRadioButton label: 'Replace part of the existing name'.
	randomNamesRadioBtn associatedRadioButtons: { addPrefixRadioBtn  . addSufixRadioBtn . 
	                                              removePartRadioBtn . replaceRadioBtn }.
	renameTypeBtns  := Dictionary newFrom: { (#randomNamesRadioBtn -> randomNamesRadioBtn ) . 
	                                   (#addPrefixRadioBtn -> addPrefixRadioBtn ) .
												 (#addSufixRadioBtn -> addSufixRadioBtn ) . 
												 (#removePartRadioBtn -> removePartRadioBtn ) .
												 (#replaceRadioBtn -> replaceRadioBtn ) }. 
	
	prefixInput  := self newTextInput autoAccept: true .
	prefixInput  placeholder: 'Enter prefix ' .
	sufixInput := self newTextInput autoAccept: true .		
	sufixInput placeholder: 'Enter sufix ' .			
	removePartInput := self newTextInput autoAccept: true .
	removePartInput placeholder: 'Enter part of the name to be removed' .
	replaceInput := self newTextInput autoAccept: true .
	replaceInput placeholder: 'Enter part of the name to be replaced' .
	replaceWithInput := self newTextInput autoAccept: true .
	replaceWithInput placeholder: 'Replace with' .		

	
	submitBtn := self newButton label: 'Rename' .
	submitBtn action: [ self submitRename ] .
]

{ #category : #initialization }
ArfRenamePresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Rename';
		initialExtent: 300@520;
		centered
]

{ #category : #'as yet unclassified' }
ArfRenamePresenter >> submitRename [
	
	^ self inform: 'Submit rename'
]
