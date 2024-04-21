#version 410 core
layout (location = 0) in vec3 pos;
layout (location = 1) in vec3 colour;

uniform mat4 projection;

out vec3 Colour;
void main() {
	gl_Position = projection * vec4(pos,1.0f);
	Colour = colour;
}
