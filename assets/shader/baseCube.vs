#version 460 core
layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;

layout(location = 2) in vec3 translate;
layout(location = 3) in vec2 handle;

uniform mat4 projection;
uniform mat4 view;

out vec2 texCoord;
out flat uvec2 fSampler;

void main() {
  // column-major
  mat4 model = mat4(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, translate.x,
                    translate.y, translate.z, 1);
  gl_Position = projection * view * model * vec4(pos, 1.0f);
  texCoord = uv;
  fSampler = uvec2(handle);
}
