#version 460 core

in vec4 fCentreColour;
in vec4 fBorderColour;
in float fBorderThickness;

in vec2 fCentre;
in float fRadius;

// ssCoord
layout(origin_upper_left) in vec4 gl_FragCoord;
out vec4 FragColour;

void main() {
  vec2 p = gl_FragCoord.xy - fCentre;
  float dist = length(p);

  float distFromBorderCentre = abs(dist - fRadius);
  // lerp: [0,borderThickness] -> [1,0] & [borderThickness,inf[ -> 0
  float a = 1 - (distFromBorderCentre / fBorderThickness);

  if (dist < fRadius) {
    FragColour = fCentreColour;
  } else {
    FragColour = vec4(fBorderColour.rgb, a);
  }
  // obiges branchless, (tatsächlich langsamer):
  /*   float cond = int(dist < fRadius);
    FragColour = cond * fCentreColour + (1 - cond) * vec4(fBorderColour.rgb, a);
  */
}