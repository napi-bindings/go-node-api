#!/bin/bash 

echo Start building process ...  && \

echo Cleanup previous build ... && \
rm -rf *.o && \
rm -rf *.a && \
echo Cleanup process completed. && \

# BUILD N-API STUB
echo Start building N-API stub library ...  && \
export MACOSX_DEPLOYMENT_TARGET=10.10 && \
gcc -c node_api.c && \
ar -rcs libnode_api.a node_api.o && \
ranlib libnode_api.a && \
echo N-API stub library successfully builded. && \

echo Building process successfully ended.