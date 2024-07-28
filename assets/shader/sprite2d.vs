#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 texCoord;
layout(location = 2) in vec4 offset;

uniform mat4 projection;

out vec2 fUV;

void main() {
  fUV = texCoord;
  gl_Position = projection * vec4(pos * offset.zw + offset.xy, 0.0, 1.0);
}