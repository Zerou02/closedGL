#version 460 core
out vec4 FragColor;

in vec4 fColour;
void main() { FragColor = vec4(fColour); }