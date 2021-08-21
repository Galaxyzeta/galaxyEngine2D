#version 410 core

in vec3 aPos;
in vec4 inputColor;

out vec4 vertColor;

void main()
{
    vertColor = inputColor;
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}