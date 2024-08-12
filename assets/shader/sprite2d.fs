#version 410 core 
out vec4 fragColour; 
in vec2 texCoord;

uniform sampler2D tex;
uniform vec4 colour;
void main() {
	//fragColour = mix(texture(tex,texCoord),texture(tex2,texCoord),0.2f) * vec4(color,1.0f);
	//	fragColour = texture(tex2,texCoord);
	//fragColour = mix(texture(tex,texCoord),texture(tex2,texCoord),0.2f);
	//fragColour = vec4(color,1.0f);
	//fragColour = vec4(0.0f,0.0f,1.0f,1.0f);
	vec4 tex = texture(tex,texCoord);
	fragColour = vec4(tex*colour);
}