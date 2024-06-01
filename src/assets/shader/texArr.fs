uniform sampler2DArray u_tex;
in vec2 f_uv;
out vec4 FragColor;
uniform float texArr;

in vec2 TexCoord2;
void main() {
        vec3 texCoordTest = vec3( f_uv.x , f_uv.y , texArr);
        vec4 color = texture(u_tex  , texCoordTest);            
        FragColor = color ; 
};