#version 460 core
layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;

layout(location = 2) in vec3 translate;

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
out flat uvec2 fSampler;
out uint fData;
out float fVertID;

struct Data {
  uvec2 handle;
  uint mask;
  uint mask2;
};

layout(binding = 1, std430) readonly buffer ssbo { Data values[]; };

void main() {
  // column-major
  mat4 model = mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, translate.x,
                    translate.y, translate.z, 1);
  gl_Position = projection * view * model * vec4(pos, 1.0f);
  texCoord = uv;
  fSampler = values[gl_InstanceID].handle;
  fData = values[gl_InstanceID].mask;
  fVertID = gl_VertexID;
}
