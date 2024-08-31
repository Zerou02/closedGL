#version 460 core
out vec4 FragColor;

in vec4 fColour;
in vec2 fUV;
in float fSignMultiplier;
void main() {
  // FragColor = fColour;
  float val = fUV.x * fUV.x - fUV.y; //<0 = innerhalb
  val = fSignMultiplier * val;
  if (val <= 0) {
    FragColor = vec4(1, 1, 1, 0.25);
  } else {
    FragColor = vec4(1, 0, 0, 0);
  }
  // FragColor = vec4(1, 1, 1, 1);
}