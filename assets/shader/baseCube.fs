
#version 460 core
#extension GL_ARB_bindless_texture : require

out vec4 fragColour;

in vec2 texCoord;
in flat uvec2 fSampler;

in flat float test;
void main() { fragColour = texture(sampler2D(fSampler), texCoord); }