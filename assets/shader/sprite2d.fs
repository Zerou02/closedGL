#version 430 core
#extension GL_ARB_bindless_texture : require

out vec4 FragColor;

in vec2 fUV;
in flat  uvec2 fSampler;
in flat float fDivisor;
void main() {
  sampler2D tex = sampler2D(fSampler);
  FragColor = texture(tex, fUV*fDivisor);
 //FragColor = vec4(1,1,1,1);
}