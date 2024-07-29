
#version 460 core
out vec4 fragColour;
in vec2 texCoord;

uniform sampler2D tex;

void main() { fragColour = vec4(texture(tex, texCoord)); }