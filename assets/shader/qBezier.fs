#version 410 core

out vec4 FragColour;

in vec2 fP1;
in vec2 fP2;
in vec2 fCp;
in vec2 fDimX;

layout(origin_upper_left) in vec4 gl_FragCoord;

void main() {

  float w = 100;
  float m = 400;
  // (1/m)(x-w),m=width,w=posX: lerp: dim.x -> [0,1]
  // float t = (1 / fDimX[1]) * (gl_FragCoord.x - fDimX[0]);
  float t = (1 / m) * (gl_FragCoord.x - w);

  vec2 fP1a = vec2(100, 100);
  vec2 fP2a = vec2(300, 200);
  vec2 fCpa = vec2(300, 300);
  /*   vec2 eq =
        (1 - t) * ((1 - t) * fP1a + t * fCpa) + t * ((1 - t) * fCpa + t * fP2a);
   */
  /*   fCp = vec2(600, 200);
    fP2 = vec2(300, 300); */
  vec2 r = mix(fP1a, fCpa, t);
  vec2 s = mix(fCpa, fP2a, t);
  vec2 eq = mix(r, s, t);

  // float dist = abs(eq[1] - gl_FragCoord[1]);
  float dist = distance(eq, gl_FragCoord.xy);

  if (dist < 1) {
    FragColour = vec4(1, 1, 1.0, 1 - dist);
  } else {
    FragColour = vec4(t, 0, 1.0, 1);
  }
}