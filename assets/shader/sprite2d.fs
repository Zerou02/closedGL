#version 460 core
out vec4 FragColor;

uniform sampler2D text;

in vec4 fColour;

void main() {
  vec4 sampled = texture(text, fColour.xy);
  FragColor = vec4(sampled);
}