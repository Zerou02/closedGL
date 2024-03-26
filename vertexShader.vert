#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
uniform vec2 offset;
out vec3 color;
void main() {
    gl_Position = vec4(aPos, 1.0f) + vec4(offset,0.0f,0.0f);
	color = aColor;
}
