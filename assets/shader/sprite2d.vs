#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec4 offset;
layout(location = 2) in vec4 colour;
layout(location = 3) in vec4 colour2;

uniform mat4 projection;

out vec2 texCoord;

void main() {
  if (gl_VertexID == 0) {
    texCoord = vec2(colour.x, colour2.x);
  } else if (gl_VertexID == 1) {
    texCoord = vec2(colour.y, colour2.y);
  } else if (gl_VertexID == 2) {
    texCoord = vec2(colour.z, colour2.z);
  } else if (gl_VertexID == 3) {
    texCoord = vec2(colour.w, colour2.w);
  }

  gl_Position = projection * vec4(pos * offset.zw + offset.xy, 0.0, 1.0);
}