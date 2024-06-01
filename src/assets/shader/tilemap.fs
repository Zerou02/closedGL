#version 410 core 
out vec4 fragColour; 
in vec2 texCoord;
in float outTexID;



uniform sampler2D tex0;
uniform sampler2D tex1;
uniform sampler2D tex2;
uniform sampler2D tex3;
uniform sampler2D tex4;
uniform sampler2D tex5;
uniform sampler2D tex6;
uniform sampler2D tex7;
uniform sampler2D tex8;
uniform sampler2D tex9;
uniform sampler2D tex10;
uniform sampler2D tex11;
uniform sampler2D tex12;
uniform sampler2D tex13;
uniform sampler2D tex14;
uniform sampler2D tex15; 

void main() {

 	if(outTexID == 0){
		fragColour = texture(tex0,texCoord);
	}else if(outTexID == 1){
		fragColour = texture(tex1,texCoord);
	}else if(outTexID == 2){
		fragColour = texture(tex2,texCoord);
	}
	else if(outTexID == 3){
		fragColour = texture(tex3,texCoord);
	}
	else if(outTexID == 4){
		fragColour = texture(tex4,texCoord);
	}else if(outTexID == 5){
		fragColour = texture(tex5,texCoord);
	}else if(outTexID == 6){
		fragColour = texture(tex6,texCoord);
	}else if(outTexID == 7){
		fragColour = texture(tex7,texCoord);
	}else if(outTexID == 8){
		fragColour = texture(tex8,texCoord);
	}else if(outTexID == 9){
		fragColour = texture(tex9,texCoord);
	}else if(outTexID == 10){
		fragColour = texture(tex10,texCoord);
	}else if(outTexID == 11){
		fragColour = texture(tex11,texCoord);
	}else if(outTexID == 12){
		fragColour = texture(tex12,texCoord);
	}else if(outTexID == 13){
		fragColour = texture(tex13,texCoord);
	}else if(outTexID == 14){
		fragColour = texture(tex14,texCoord);
	}else if(outTexID == 15){
		fragColour = texture(tex15,texCoord);
	}  
}