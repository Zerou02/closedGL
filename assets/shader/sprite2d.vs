#version 430 core
layout(location = 0) in vec4 pos; //baseBuffer: <posX,posY,uv>
layout(location = 1) in vec4 offset;
layout(location = 2) in vec4 uvData; // posX,posY,dimX,dimY
layout(location = 3) in vec2 cellSpriteSize;

uniform mat4 projection;

layout(binding = 1, std430) readonly buffer ssbo { uvec2 values[]; };

out vec2 fUV;
out flat uvec2 fSampler;
out flat float fDivisor;

void main() {
  fDivisor = cellSpriteSize[0]/cellSpriteSize[1];
  fUV = vec2(uvData[0]+pos.z*uvData[2],uvData[1]+pos.w*uvData[3]);
  float test = pos[0];
  fSampler = values[int(gl_InstanceID)];
  gl_Position = projection * vec4(pos.xy * offset.zw + offset.xy, 0.0, 1.0);
}