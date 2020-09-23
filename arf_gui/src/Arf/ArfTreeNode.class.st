Class {
	#name : #ArfTreeNode,
	#superclass : #Object,
	#instVars : [
		'element',
		'children'
	],
	#category : #Arf
}

{ #category : #accessing }
ArfTreeNode >> children [ 
	^ children
]

{ #category : #accessing }
ArfTreeNode >> element [
	^ element
]

{ #category : #accessing }
ArfTreeNode >> element: aElement [
	element := aElement
]

{ #category : #initialization }
ArfTreeNode >> initialize [ 
		
	element := ArfFile new initialize . 
]
