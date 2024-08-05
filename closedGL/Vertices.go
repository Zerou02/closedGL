package closedGL

var CompressedCubeVertices = []byte{
	31, 26, 8, 8, 13, 31, 5, 23, 30, 30, 12, 5, 14, 8, 1, 1, 7, 14, 19, 26, 28, 19, 28, 21, 24, 17, 3, 3, 10, 24, 2, 16, 21, 21, 7, 2,
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

var CubeVertices = cube
