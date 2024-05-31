#version 410 core
in vec2 TexCoords;
out vec4 Colour;

uniform sampler2D text;
uniform vec3 colour;

void main()
{    
    vec4 sampled = texture(text, TexCoords);
    Colour = vec4(sampled);
}