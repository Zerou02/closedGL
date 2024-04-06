#version 410 core 
out vec4 fragColour; 

uniform vec3 objectColour;
uniform vec3 lightColour;
uniform vec3 lightPos;
uniform vec3 viewPos;

in vec3 normal;
in vec3 fragWorldPos;
void main() {
    float specularStrength = 0.5;
    vec3 normalized = normalize(normal);
    vec3 lightDir = normalize(lightPos-fragWorldPos);
    vec3 viewDir = normalize(viewPos-fragWorldPos);
    vec3 reflectDir = reflect(-lightDir,normalized);
    //verhindert negative Zahlen
    float diff = max(dot(normalized,lightDir),0.0);
    vec3 diffuse = diff * lightColour;
    int shininess = 32;
    float spec = pow(max(dot(viewDir,reflectDir),0.0),shininess);
    vec3 specular = specularStrength * spec * lightColour;
    float ambientStrength = 0.1;
    vec3 ambient = ambientStrength *lightColour;
    vec3 res = (ambient+diffuse+specular)* objectColour;
	fragColour = vec4(res,1.0);
}