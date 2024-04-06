#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;

uniform vec2 offset;
uniform mat4 transform;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
out vec3 color;
out vec2 texCoord;
void main() {
	gl_Position = projection * view * model * vec4(aPos,1.0f);
       //gl_Position = vec4(aPos, 1.0f) + vec4(offset,0.0f,0.0f);
	color = aColor;
	texCoord = aTexCoord;
}
