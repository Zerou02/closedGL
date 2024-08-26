#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec4 offset;
layout(location = 2) in vec4 colour;

uniform mat4 projection;

out vec4 fColour;

void main() {
  fColour = colour;
  gl_Position = projection * vec4(pos * offset.zw + offset.xy, 0.0, 1.0);
}