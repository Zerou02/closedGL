#version 410 core 
out vec4 fragColour; 
in vec2 texCoord;

uniform sampler2D tex;
uniform vec4 colour;
uniform float c;
void main() {
	vec4 tex = texture(tex,texCoord*32/1024);
	fragColour = vec4(tex.rgb,1);
}