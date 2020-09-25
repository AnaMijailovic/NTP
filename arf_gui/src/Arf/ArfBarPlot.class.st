"
Inherits RSBarPlot class, overrides its renderIn: method
and enables showing x labels.


"
Class {
	#name : #ArfBarPlot,
	#superclass : #RSBarPlot,
	#instVars : [
		'xLabels'
	],
	#category : #'Arf-presenters'
}

{ #category : #rendering }
ArfBarPlot >> renderIn: canvas [
	| color width myScale label composite| 
	
	self checkAssertion.
	color := self computeColor.
   "myScale := self barScale."
	myScale := xScale.
	width := self barWidth.
	width ifNil: [ width := myScale rangeBand ].
	xValues doWithIndex: [ :xt : i | 
		| yt origin corner r |
		yt := yValues at: i.
		origin := (myScale scale: xt) @ (yScale scale: yt).
		corner := origin x @ (yScale scale: 0).
		((self isPointWellDefined: origin) and: [ self isPointWellDefined: corner ]) ifTrue: [
			r := Rectangle origin: origin corner: corner + (width@0).
			
			"Show x label"
			label := RSLabel new text: (xLabels at: i) ; fontSize: 5.
	    
			composite := RSComposite new.
			composite add: label; 
						 position: corner + ((width/2)@5).
			canvas add: composite .
		   canvas add: (self shape copy
				color: color;
				fromRectangle:  r;
				yourself).
				
			  ] ].
]

{ #category : #accessing }
ArfBarPlot >> xLabels: aYLabels [
	
	xLabels := aYLabels
]
