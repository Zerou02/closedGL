#version 460 core

#extension GL_ARB_bindless_texture : require
#extension GL_NV_gpu_shader5 : require

out vec4 FragColor;

uniform sampler2D text;

layout(binding = 1, std430) readonly buffer ssbo { sampler2D values[]; };

in vec2 fUV;

void main() {
  sampler2D tex = values[0];
  vec4 sampled = texture(tex, fUV);
  FragColor = vec4(sampled);
}