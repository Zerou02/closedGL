#version 460 core

layout(location = 0) in uint instData;
layout(location = 2) in vec3 translate;
layout(location = 3) in vec2 handle;

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
out flat uvec2 fSampler;

const float vertexPos[2] = float[](-0.5, 0.5);

layout(binding = 0, std430) readonly buffer ssbo { uint testData[]; };

out flat float test;
void main() {
  test = 0;
  if (instData == 5) {
    test = 1;
  }
  test = 0;
  // x,y,z,u,v
  int data = int(testData[(instData & 7) * 6 + gl_VertexID]);
  int v = (data >> 0) & 1;
  int u = (data >> 1) & 1;
  int z = (data >> 2) & 1;
  int y = (data >> 3) & 1;
  int x = (data >> 4) & 1;
  // column-major
  mat4 model = mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, translate.x,
                    translate.y, translate.z, 1);
  gl_Position = projection * view * model *
                vec4(vec3(vertexPos[x], vertexPos[y], vertexPos[z]), 1.0f);
  texCoord = vec2(u, v);
  fSampler = uvec2(handle);
}
