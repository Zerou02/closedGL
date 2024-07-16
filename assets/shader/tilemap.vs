#version 410 core
layout (location = 0) in vec2 pos;
layout (location = 1) in vec2 texC;
layout (location = 2) in float texID;

uniform mat4 model;
uniform mat4 projection;

out vec2 texCoord;
out float outTexID;
void main() {
	gl_Position = projection * model * vec4(pos,0,1.0f);

	texCoord = texC;
    outTexID = texID;
}
