#version 410 core
in vec4 Colour;
out vec4 FragColour;

void main()
{    
    FragColour = vec4(Colour);
}