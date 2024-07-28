#version 460 core

#extension GL_ARB_bindless_texture : require

out vec4 FragColor;

layout(binding = 1, std430) readonly buffer ssbo { sampler2D values[]; };

in vec2 fUV;

void main() {
  sampler2D tex = values[0];
  vec4 sampled = texture(tex, fUV);
  FragColor = vec4(sampled);
}