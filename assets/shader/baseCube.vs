#version 460 core
layout(location = 0) in vec2 handle;
layout(location = 1) in vec3 pos;
layout(location = 2) in vec2 uv;

layout(location = 3) in vec3 translate;

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
out flat vec2 fSampler;
out flat uvec2 fSamplerSSBO;

out flat float test;
struct Data {
  uvec2 handle;
};

layout(binding = 1, std430) readonly buffer ssbo { Data values[]; };

void main() {
  // column-major
  mat4 model = mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, translate.x,
                    translate.y, translate.z, 1);
  gl_Position = projection * view * model * vec4(pos, 1.0f);
  texCoord = uv;
  fSampler = handle;
  fSamplerSSBO = values[0].handle;
  test = 0;
  if (uvec2(fSampler) == fSamplerSSBO) {
    test = 1;
  }
}
