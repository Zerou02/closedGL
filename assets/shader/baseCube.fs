
#version 460 core
#extension GL_ARB_bindless_texture : require

out vec4 fragColour;

in vec2 texCoord;
in float fInstanceID;

layout(binding = 1, std430) readonly buffer ssbo { sampler2D values[]; };

void main() {
  sampler2D tex = values[int(fInstanceID)];
  fragColour = texture(tex, texCoord);
  // fragColour = vec4(1, 1, 1, 1);
}