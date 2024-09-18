#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 uv;
layout(location = 2) in float signMultiplier;

uniform mat4 projection;
uniform mat4 view;

out vec4 fColour;
out vec2 fUV;
out float fSignMultiplier;
void main() {

  fSignMultiplier = signMultiplier;
  fUV = uv;
  fColour = vec4(1, 1, 1, 1);
  gl_Position = projection * vec4(pos, 0.0, 1.0);
}