#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec4 colour;

uniform mat4 projection;
uniform mat4 view;

out vec4 fColour;

void main() {
  fColour = colour;
  gl_Position = projection * view * vec4(pos, 0.0, 1.0);
}