#version 410 core

// 1.)4texID,5u,5v,5x,5y,5z,3side
// 2.)18xyz
layout(location = 0) in uvec2 instData;

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
out flat uvec2 fSampler;

const float vertexPos[2] = float[](0, 1);

layout(binding = 0, std430) readonly buffer baseMeshSSBO {
  uint baseMeshData[36];
  uvec2 handles[];
};
layout(binding = 1, std430) readonly buffer meshSSBO { uint meshData[]; };

void main() {

  // x,y,z,u,v
  int data = int(baseMeshData[(instData.x & 7) * 6 + gl_VertexID]);
  int v = (data >> 0) & 1;
  int u = (data >> 1) & 1;
  int z = (data >> 2) & 1;
  int y = (data >> 3) & 1;
  int x = (data >> 4) & 1;
  // column-major
  mat4 model = mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0,
                    float(((instData.x >> 13) & 31) + meshData[0]),
                    float(((instData.x >> 8) & 31) + meshData[1]),
                    float(((instData.x >> 3) & 31) + meshData[2]), 1);

  uint sizeZ = instData.y & 63;
  uint sizeY = (instData.y >> 6) & 63;
  uint sizeX = (instData.y >> 12) & 63;
  gl_Position = projection * view * model *
                vec4(vec3(vertexPos[x], vertexPos[y], vertexPos[z]) *
                         vec3(sizeX, sizeY, sizeZ),
                     1.0f);

  uint texU = instData.x >> 23 & 31;
  uint texV = instData.x >> 18 & 31;

  texCoord = vec2(texU + u, texV + v);
  fSampler = handles[instData.x >> 27 & 15];
}
