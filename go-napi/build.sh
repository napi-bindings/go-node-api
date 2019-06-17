echo Cleaning previous build ... && \
rm -rf *.a && \
rm -rf libgoaddon.h && \
echo Start building ... && \
# Remember for Node.js version less than 12 the MACOSX_DEPLOYMENT_TARGET need to 
# be set to 10.7
export MACOSX_DEPLOYMENT_TARGET=10.10 && go build -a -x -o libgoaddon.a -buildmode=c-archive . && \
echo Build finished.

