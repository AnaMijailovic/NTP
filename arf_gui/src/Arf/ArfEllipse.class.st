Class {
	#name : #ArfEllipse,
	#superclass : #RSEllipse,
	#traits : 'SpTContextMenu',
	#classTraits : 'SpTContextMenu classTrait',
	#instVars : [
		'#selectionMode',
		'#activationBlock',
		'#activateOnSingleClick',
		'#model'
	],
	#category : #Arf
}

{ #category : #accessing }
ArfEllipse >> activateOnSingleClick [
	
	activateOnSingleClick := true
]

{ #category : #testing }
ArfEllipse >> activatesOnSingleClick [
	
	^ activateOnSingleClick
]

{ #category : #accessing }
ArfEllipse >> contextMenu: aValuable [
	"Set the block used to define a context menu for the presenter"
	self inform: 'context menu'.
	contextMenu := aValuable
]

{ #category : #initialization }
ArfEllipse >> initialize [
	
	self class initializeSlots: self.
	super initialize.
]
