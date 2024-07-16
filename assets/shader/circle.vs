#version 460 core
layout(location = 0) in vec2 pos;
layout(location = 1) in vec4 data; // posX,posY,r,borderThickness
layout(location = 2) in vec4 centreColour;
layout(location = 3) in vec4 borderColour;

uniform mat4 projection;

out vec4 fCentreColour;
out vec4 fBorderColour;
out float fBorderThickness;

out vec2 fCentre;
out float fRadius;

void main() {
  fCentreColour = centreColour;
  fBorderColour = borderColour;
  fBorderThickness = data.w;
  fRadius = data.z / 2;
  fCentre = vec2(data.x + data.z / 2, data.y + data.z / 2);

  gl_Position = projection * vec4(pos * data.z + data.xy, 0.0, 1.0);
}
