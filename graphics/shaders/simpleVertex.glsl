#version 410 core

in vec3 aPos;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main()
{
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}