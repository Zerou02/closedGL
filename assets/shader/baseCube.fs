
#version 460 core
#extension GL_ARB_bindless_texture : require

out vec4 fragColour;

in vec2 texCoord;
in flat uvec2 fSampler;

in flat float test;

void main() { // fragColour = vec4(0, 1, 0, 1); }
  if (test == 1) {
    fragColour = vec4(1, 1, 1, 1);
  } else {
    fragColour = texture(sampler2D(fSampler), texCoord);
  }
}