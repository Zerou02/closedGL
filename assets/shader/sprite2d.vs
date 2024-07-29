#version 460 core
layout(location = 0) in vec4 pos; //<posX,posY,uvX,uvY>
layout(location = 1) in vec4 offset;

uniform mat4 projection;

out vec2 fUV;
out float fInstanceID;
void main() {
  fUV = pos.zw;
  fInstanceID = gl_InstanceID;
  gl_Position = projection * vec4(pos.xy * offset.zw + offset.xy, 0.0, 1.0);
}