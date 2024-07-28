#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec4 offset;
layout(location = 2) in vec4 colour;
layout(location = 3) in vec4 colour2;

uniform mat4 projection;

out vec4 fColour;

void main() {
  fColour = colour;
  if (gl_VertexID == 0) {
    fColour = vec4(colour.x, colour2.x, 0, 0);
  } else if (gl_VertexID == 1) {
    fColour = vec4(colour.y, colour2.y, 0, 0);
  } else if (gl_VertexID == 2) {
    fColour = vec4(colour.z, colour2.z, 0, 0);
  } else if (gl_VertexID == 3) {
    fColour = vec4(colour.w, colour2.w, 0, 0);
  }
  gl_Position = projection * vec4(pos * offset.zw + offset.xy, 0.0, 1.0);
}