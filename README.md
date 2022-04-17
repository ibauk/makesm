# makesm

This utility makes working installations from the ScoreMaster development environment.

**makesm** **-target** *path* [**-ok**] [**-rblr**] [**-sm3** *path*] [**-runsm** *path*] [**-smpatch** *path*] [**-ebcfetch** *path*] [**-php** *path*]

**target**
: specifies the path into which an installation will be built

**ok**
: forces overwrite of existing installation

**rblr**
: includes extra certificates and images for use with the RBLR1000 event

The other *path* specs point to working copies of ScoreMaster code, utilities and PHP. The defaults all use Windows paths.
