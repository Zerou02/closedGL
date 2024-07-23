#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec4 offset;
layout(location = 2) in vec2 p1;
layout(location = 3) in vec2 p2;
layout(location = 4) in vec2 cp;

uniform mat4 projection;

out vec2 fP1;
out vec2 fP2;
out vec2 fCp;
out vec2 fDimX;
void main() {
  fDimX = vec2(offset.x, offset.z);
  fP1 = p1;
  fP2 = p2;
  fCp = cp;

  gl_Position = projection * vec4(pos * offset.zw + offset.xy, 0.0, 1.0);
}