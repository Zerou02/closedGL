#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
//Dient dazu, Model-Matrix zu konstruieren
layout (location = 2) in vec3 worldPos; 

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
void main() {
	//translation-matrix; column-major; 4x4homogenous
	mat4 modelMat = mat4(
		1,0,0,0,
		0,1,0,0,
		0,0,1,0,
		worldPos.x,worldPos.y,worldPos.z,1
	);
	gl_Position = projection * view * modelMat * vec4(aPos,1.0f);
	texCoord = aTexCoord;
}
