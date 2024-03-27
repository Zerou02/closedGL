#version 410 core 
out vec4 fragColour; 
in vec3 color; 
in vec2 texCoord;

uniform sampler2D tex;
uniform sampler2D tex2;	
void main() {
		fragColour = mix(texture(tex,texCoord),texture(tex2,texCoord),0.2f) * vec4(color,1.0f);
		//fragColour = mix(texture(tex,texCoord),texture(tex2,texCoord),0.2f);
		//fragColour = vec4(color,1.0f);
		//fragColour = vec4(1.0f,1.0f,1.0f,1.0f);

}

