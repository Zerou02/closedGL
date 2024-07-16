package closedGL

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

var indicesQuad = []uint32{
	0, 1, 3,
	1, 2, 3,
}

var fullQuad = []float32{
	//pos;;tex
	1, 0, 0, 1.0, 0.0, // top right
	0, 0, 0, 0.0, 0.0, // top left
	1, 1, 0, 1.0, 1.0, // bottom right
	1, 1, 0, 1.0, 1.0, // bottom right
	0, 0, 0, 0.0, 0.0, // top left
	0, 1, 0, 0.0, 1.0, // bottom left
}
