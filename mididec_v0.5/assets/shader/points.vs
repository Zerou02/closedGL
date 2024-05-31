#version 410 core
layout (location = 0) in vec2 pos;
layout (location = 1) in vec4 colour;

uniform mat4 projection;

out vec4 Colour;
void main() {
	gl_Position = projection * vec4(pos.xy,0,1.0f);
	Colour = colour;
}
