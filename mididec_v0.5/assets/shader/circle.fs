#version 410 core
in vec4 Colour;
//ssCoord
layout (origin_upper_left) in vec4 gl_FragCoord;
out vec4 FragColour;

uniform vec2 centre;
uniform float radius;
uniform float borderThickness;
uniform vec4 borderColour;
uniform vec4 centreColour;
void main()
{
    vec2 p = gl_FragCoord.xy-centre;
    float dist = length(p);

    float distFromBorderCentre = abs(dist-radius);
    //lerp: [0,borderThickness] -> [1,0] ^ [borderThickness,inf[ -> 0
    float a = 1-(1/borderThickness*distFromBorderCentre);
    if(dist < radius){
        FragColour = vec4(centreColour);
    }else{
        FragColour = vec4(borderColour.rgb,a);
    }
} 