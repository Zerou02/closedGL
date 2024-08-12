#version 410 core
#extension GL_ARB_bindless_texture : require

out vec4 FragColor;

layout(binding = 1, std430) readonly buffer ssbo { sampler2D values[]; };

in vec2 fUV;
in float fInstanceID;

void main() {
  sampler2D tex = values[int(fInstanceID)];
  FragColor = texture(tex, fUV);
}