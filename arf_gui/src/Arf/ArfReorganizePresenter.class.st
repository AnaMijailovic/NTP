Class {
	#name : #ArfReorganizePresenter,
	#superclass : #SpPresenterWithModel,
	#instVars : [
		'pathInput',
		'recursiveCheckBox',
		'fileTypeRadioBtn',
		'creationRadioBtn',
		'creationDateStepInput',
		'fileSizeRadioBtn',
		'sizeStepInput',
		'reorganizeTypeBtns',
		'submitBtn'
	],
	#category : #Arf
}

{ #category : #specs }
ArfReorganizePresenter class >> defaultLayout [

		^ SpBoxLayout   newVertical 
			add: 'Path: ' withConstraints: [ :c | c height: 25];
			add: #pathInput withConstraints: [ :c | c height: 25];
			add: #recursiveCheckBox withConstraints: [ :c | c height: 45];
			add: #fileTypeRadioBtn withConstraints: [ :c | c height: 45];
			add: #creationRadioBtn withConstraints: [ :c | c height: 45];
			add: #creationDateStepInput withConstraints: [ :c | c height: 45];
			add: #fileSizeRadioBtn withConstraints: [ :c | c height: 45];
			add: #sizeStepInput;
			add: #submitBtn withConstraints: [ :c | c height: 40];
		   yourself.
]

{ #category : #specs }
ArfReorganizePresenter class >> defaultSpec [

	^ self defaultLayout 
]

{ #category : #initialization }
ArfReorganizePresenter >> initializePresenters [
	pathInput := self newTextInput autoAccept: true.
	
	recursiveCheckBox := self newCheckBox label: 'Recursive'.
	
	fileTypeRadioBtn := self newRadioButton label: 'File type'.
	creationRadioBtn := self newRadioButton label: 'Creation date'.
	fileSizeRadioBtn := self newRadioButton label: 'File size'.
	fileTypeRadioBtn associatedRadioButtons: { creationRadioBtn . fileSizeRadioBtn }.
	reorganizeTypeBtns  := Dictionary newFrom: { (#fileTypeRadioBtn -> fileTypeRadioBtn  ) . 
	                                   (#creationRadioBtn -> creationRadioBtn  ) .
												 (#fileSizeRadioBtn -> fileSizeRadioBtn  ) }. 
	
	creationDateStepInput := self newNumberInput autoAccept: true.
	creationDateStepInput placeholder: 'Enter step (in months)'; beInteger .
	sizeStepInput := self newNumberInput autoAccept: true.		
	sizeStepInput placeholder: 'Enter step'; beInteger .							

	
	submitBtn := self newButton label: 'Reorgnize'.
	submitBtn action: [ self submitReorganize ].
]

{ #category : #initialization }
ArfReorganizePresenter >> initializeWindow: aWindowPresenter [
	aWindowPresenter
		title: 'Reorganize';
		initialExtent: 300@430;
		centered
]

{ #category : #'as yet unclassified' }
ArfReorganizePresenter >> submitReorganize [
	^ self inform: 'Submit reorganize'
]
