#version 410 core
layout (location = 0) in uint vertex;

uniform mat4 projection;
uniform mat4 view;
uniform vec3 chunkOrigin;

out vec2 texCoord;

void main() {
	uint val = uint(vertex);
	uint modelZ = val & 31;
	val = val >> 5;
	uint modelY = val & 31;
	val = val >> 5;
	uint modelX = val & 31;
	val = val >> 5;
	uint texY = val & 31;
	val = val >> 5;
	uint texX = val & 31;
	val = val >> 5;
	uint ndcIDx = val & 7;

	//column-major;
	mat4 scaleMat = mat4(
		1,0,0,0,
		0,1,0,0,
		0,0,1,0,
		0,0,0,1
	);

 	mat4 modelMat = mat4(
		1,0,0,0,
		0,1,0,0,
		0,0,1,0,
		float(modelX+chunkOrigin.x),
		float(modelY+chunkOrigin.y),
		float(modelZ+chunkOrigin.z),
		1
	);
	uint x = (ndcIDx>>2) & 1;
	uint y = (ndcIDx>>1) & 1;
	uint z = (ndcIDx>>0) & 1;
	
	/* 
	float C = 1.0; 
	float far = 2000.0;  */
	gl_Position = projection * view * scaleMat* modelMat * vec4(float(x),float(y),float(z),1.0f);
	/* 	gl_Position.z = 2.0*log(gl_Position.w*C + 1)/log(far*C + 1) - 1;
    gl_Position.z *= gl_Position.w; */
	texCoord = vec2(float(texX),float(texY));
}
