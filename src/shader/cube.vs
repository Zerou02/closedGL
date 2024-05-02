#version 410 core
layout (location = 0) in float aPos;
layout (location = 1) in vec2 aTexCoord;
//Dient dazu, Model-Matrix zu konstruieren
layout (location = 2) in float worldPos; 

uniform mat4 projection;
uniform mat4 view;
uniform vec3 chunkOrigin;

out vec2 texCoord;
void main() {
	//translation-matrix; column-major; 4x4homogenous
	int worldX = (int(worldPos)>>10) & 31;
	int worldY = (int(worldPos)>>5) &  31;
	int worldZ = (int(worldPos)>>0) &  31;

 	mat4 modelMat = mat4(
		1,0,0,0,
		0,1,0,0,
		0,0,1,0,
		float(worldX)+chunkOrigin.x,
		float(worldY)+chunkOrigin.y,
		float(worldZ)+chunkOrigin.z,
		1
	); 
	int x = (int(aPos)>>2) & 1;
	int y = (int(aPos)>>1) & 1;
	int z = (int(aPos)>>0) & 1;

	gl_Position = projection * view * modelMat * vec4(float(x),float(y),float(z),1.0f);
	texCoord = aTexCoord;
}
