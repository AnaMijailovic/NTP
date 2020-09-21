Class {
	#name : #ArfInformationDialog,
	#superclass : #SpPresenter,
	#instVars : [
		'message',
		'okBtn'
	],
	#category : #Arf
}

{ #category : #initialization }
ArfInformationDialog class >> defaultLayout [
	
	^ SpBoxLayout  newVertical
			add: '' withConstraints: [ :c | c height: 25];
			add: #message withConstraints: [ :c | c height: 185];
			add: #okBtn "expand: true fill: false padding: 10";
			yourself. 
]

{ #category : #specs }
ArfInformationDialog class >> defaultSpec [

	^ self defaultLayout
]

{ #category : #initialization }
ArfInformationDialog >> initializePresenters [
	
	message := self newLabel label: ''.
	okBtn := self newButton label: 'OK';
				action: [ self window close ] 
	

]

{ #category : #initialization }
ArfInformationDialog >> initializeWindow: aWindowPresenter [
	
	aWindowPresenter 
		title: 'Information Dialog';
		initialExtent: 300@300;
		centered
]

{ #category : #accessing }
ArfInformationDialog >> message: aMessage [

	message label: aMessage
]
