
#version 460 core
#extension GL_ARB_bindless_texture : require

out vec4 fragColour;

in vec2 texCoord;
in flat vec2 fSampler;
in flat uvec2 fSamplerSSBO;

in flat float test;
void main() {
  if (test == 1) {
    fragColour = vec4(0, 1, 0, 1);
  } else {
    fragColour = vec4(1, 0, 0, 1);
  }

  fragColour = texture(sampler2D(uvec2(fSampler)), texCoord);
  //  fragColour = vec4(1, fSampler.x, fSampler.y, 1);
}