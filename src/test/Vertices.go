package closed_gl

var doubleTriangle = []float32{
	//pos ;; col;;tex
	-1.0, -1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.5, 0.5,
	-1.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0,
	//pos ;; col;;tex
	1.0, -1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.5, 0.5,
	1.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0,
}

var cube = []float32{
	//pos x,y,z;;tex s,t

	//oben
	0.5, 0.5, 0.5, 1.0, 1.0, //vorne rechts oben
	0.5, 0.5, -0.5, 1.0, 0.0, //hinten rechts oben
	-0.5, 0.5, -0.5, 0.0, 0.0, //hinten links oben
	-0.5, 0.5, -0.5, 0.0, 0.0, //hinten links oben
	-0.5, 0.5, 0.5, 0.0, 1.0, //vorne links oben
	0.5, 0.5, 0.5, 1.0, 1.0, //vorne rechts oben

	//vorne
	-0.5, -0.5, 0.5, 0.0, 1.0, //vorne links unten
	0.5, -0.5, 0.5, 1.0, 1.0, //vorne rechts unten
	0.5, 0.5, 0.5, 1.0, 0.0, //vorne rechts oben
	0.5, 0.5, 0.5, 1.0, 0.0, //vorne rechts oben
	-0.5, 0.5, 0.5, 0.0, 0.0, //vorne links oben
	-0.5, -0.5, 0.5, 0.0, 1.0, //vorne links unten

	//links
	-0.5, 0.5, 0.5, 1.0, 0.0, //vorne links oben
	-0.5, 0.5, -0.5, 0.0, 0.0, //hinten links oben
	-0.5, -0.5, -0.5, 0.0, 1.0, //hinten links unten
	-0.5, -0.5, -0.5, 0.0, 1.0, //hinten links unten
	-0.5, -0.5, 0.5, 1.0, 1.0, //vorne links unten
	-0.5, 0.5, 0.5, 1.0, 0.0, //vorne links oben

	// rechts
	0.5, -0.5, -0.5, 1.0, 1.0, //hinten rechts unten
	0.5, 0.5, -0.5, 1.0, 0.0, //hinten rechts oben
	0.5, 0.5, 0.5, 0.0, 0.0, //vorne rechts oben
	0.5, -0.5, -0.5, 1.0, 1.0, //hinten rechts unten
	0.5, 0.5, 0.5, 0.0, 0.0, //vorne rechts oben
	0.5, -0.5, 0.5, 0.0, 1.0, //vorne rechts unten

	//hinten
	0.5, 0.5, -0.5, 0.0, 0.0, //hinten rechts oben
	0.5, -0.5, -0.5, 0.0, 1.0, //hinten rechts unten
	-0.5, -0.5, -0.5, 1.0, 1.0, //hinten links unten
	-0.5, -0.5, -0.5, 1.0, 1.0, //hinten links unten
	-0.5, 0.5, -0.5, 1.0, 0.0, //hinten links oben
	0.5, 0.5, -0.5, 0.0, 0.0, //hinten rechts oben

	//unten
	-0.5, -0.5, -0.5, 1.0, 0.0, //links unten hinten
	0.5, -0.5, -0.5, 0.0, 0.0, //hinten rechts unten
	0.5, -0.5, 0.5, 0.0, 1.0, //vorne rechts unten
	0.5, -0.5, 0.5, 0.0, 1.0, //vorne rechts unten
	-0.5, -0.5, 0.5, 1.0, 1.0, //vorne links unten
	-0.5, -0.5, -0.5, 1.0, 0.0, //hinten links unten
}

var cubeVertices = []float32{
	//pos x,y,z;;tex s,t
	//oben
	1, 1, 1, 1.0, 1.0, //vorne rechts oben
	1, 1, 0.0, 1.0, 0.0, //hinten rechts oben
	0.0, 1, 0.0, 0.0, 0.0, //hinten links oben
	0.0, 1, 0.0, 0.0, 0.0, //hinten links oben
	0.0, 1, 1, 0.0, 1.0, //vorne links oben
	1, 1, 1, 1.0, 1.0, //vorne rechts oben

	//vorne
	0.0, 0.0, 1, 0.0, 1.0, //vorne links unten
	1, 0.0, 1, 1.0, 1.0, //vorne rechts unten
	1, 1, 1, 1.0, 0.0, //vorne rechts oben
	1, 1, 1, 1.0, 0.0, //vorne rechts oben
	0.0, 1, 1, 0.0, 0.0, //vorne links oben
	0.0, 0.0, 1, 0.0, 1.0, //vorne links unten

	//links
	0.0, 1, 1, 1.0, 0.0, //vorne links oben
	0.0, 1, 0.0, 0.0, 0.0, //hinten links oben
	0.0, 0.0, 0.0, 0.0, 1.0, //hinten links unten
	0.0, 0.0, 0.0, 0.0, 1.0, //hinten links unten
	0.0, 0.0, 1, 1.0, 1.0, //vorne links unten
	0.0, 1, 1, 1.0, 0.0, //vorne links oben

	//hinten
	1, 1, 0.0, 0.0, 0.0, //hinten rechts oben
	1, 0.0, 0.0, 0.0, 1.0, //hinten rechts unten
	0.0, 0.0, 0.0, 1.0, 1.0, //hinten links unten
	0.0, 0.0, 0.0, 1.0, 1.0, //hinten links unten
	0.0, 1, 0.0, 1.0, 0.0, //hinten links oben
	1, 1, 0.0, 0.0, 0.0, //hinten rechts oben
	//rechts
	1, 0.0, 0.0, 1.0, 1.0, //hinten rechts unten
	1, 1, 0.0, 1.0, 0.0, //hinten rechts oben
	1, 1, 1, 0.0, 0.0, //vorne rechts oben
	1, 0.0, 0.0, 1.0, 1.0, //hinten rechts unten
	1, 1, 1, 0.0, 0.0, //vorne rechts oben
	1, 0.0, 1, 0.0, 1.0, //vorne rechts unten

	//unten
	0.0, 0.0, 0.0, 1.0, 0.0, //links unten hinten
	1, 0.0, 0.0, 0.0, 0.0, //hinten rechts unten
	1, 0.0, 1, 0.0, 1.0, //vorne rechts unten
	1, 0.0, 1, 0.0, 1.0, //vorne rechts unten
	0.0, 0.0, 1, 1.0, 1.0, //vorne links unten
	0.0, 0.0, 0.0, 1.0, 0.0, //hinten links unten

}

var quad = []float32{
	//pos;;col;;tex
	1, 0, 0, 0.5, 0.0, // top right
	1, 1, 0, 0.5, 1.0, // bottom right
	0, 1, 0, 0.0, 1.0, // bottom left
	0, 0, 0, 0.0, 0.0, // top left
}

var indicesQuad = []uint32{
	0, 1, 3,
	1, 2, 3,
}

var fullQuad = []float32{
	//pos;;tex
	1, 0, 0, 1.0, 0.0, // top right
	1, 1, 0, 1.0, 1.0, // bottom right
	0, 1, 0, 0.0, 1.0, // bottom left
	0, 0, 0, 0.0, 0.0, // top left
}

var cube8 = []float32{
	//pos;;col;;tex
	-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0, //0
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, //1
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0, //2
	-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 0.0, //3
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0, //4
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, //5
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, //6
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 1.0, //7
}

var cube24 = []float32{
	//pos;;col;;tex
	//front: 0,1,7,5
	-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0, 0.5, //0
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.5, 0.5, //1
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.5, 0, //2
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0, //3
	//right: 1,2,6,7
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0, 0.5, //4
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.5, 0.5, //5
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.5, 0, //6
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0, //7
	//back: 2,3,4,6
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0, 0.5, //8
	-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.5, 0.5, // 9
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.5, 0, //10
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0, //11
	//left: 3,0,5,4
	-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0, 0.5, //12
	-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.5, 0.5, //13
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.5, 0, //14
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0, //15
	//top: 5,7,6,4
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.5, 0.5, //16
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1, 0.5, //17
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1, 0.0, //18
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.5, 0, //19
	//bottom: 0,1,2,3
	-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 1, //20
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.5, 1.0, //21
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.5, 0.5, //22
	-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.5, //23
}

var indicesCube24 = []uint32{
	//bo
	20, 21, 22,
	22, 23, 20,
	//f
	0, 1, 2,
	2, 3, 0,
	//r
	4, 5, 6,
	6, 7, 4,
	//ba
	8, 9, 10,
	10, 11, 8,
	//l
	12, 13, 14,
	14, 15, 12,
	//top
	16, 17, 18,
	18, 19, 16,
}
