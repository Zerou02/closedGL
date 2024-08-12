#version 410 core
#extension GL_ARB_bindless_texture : require

out vec4 fragColour;

in vec2 texCoord;
in flat uvec2 fSampler;

void main() {
  fragColour = texture(sampler2D(fSampler), texCoord * 32 / 1024);
  // fragColour = vec4(1, 1, 1, 1);
}