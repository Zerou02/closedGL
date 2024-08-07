#version 460 core

// 4texID,5u,5v,5x,5y,5z,3side
layout(location = 0) in uint instData;

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
out flat uvec2 fSampler;

const float vertexPos[2] = float[](-0.5, 0.5);

layout(binding = 0, std430) readonly buffer baseMeshSSBO {
  uint baseMeshData[36];
  uvec2 handles[];
};
layout(binding = 1, std430) readonly buffer meshSSBO { uint meshData[]; };

void main() {

  // x,y,z,u,v
  int data = int(baseMeshData[(instData & 7) * 6 + gl_VertexID]);
  int v = (data >> 0) & 1;
  int u = (data >> 1) & 1;
  int z = (data >> 2) & 1;
  int y = (data >> 3) & 1;
  int x = (data >> 4) & 1;
  // column-major
  mat4 model = mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0,
                    float(((instData >> 13) & 31) + meshData[0]),
                    float(((instData >> 8) & 31) + meshData[1]),
                    float(((instData >> 3) & 31) + meshData[2]), 1);
  gl_Position =
      projection * view * model *
      vec4(vec3(vertexPos[x], vertexPos[y], vertexPos[z]) * vec3(1, 1, 1),
           1.0f);

  uint texU = instData >> 23 & 31;
  uint texV = instData >> 18 & 31;

  texCoord = vec2(texU + u, texV + v);
  fSampler = handles[instData >> 27 & 15];
}
