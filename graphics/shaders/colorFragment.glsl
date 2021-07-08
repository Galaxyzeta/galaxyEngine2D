#version 410 core
uniform sampler2D tex;

in vec4 vertColor;

void main() {
    gl_FragColor = vertColor;
}