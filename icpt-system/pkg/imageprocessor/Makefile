# Makefile for C++ Image Processor Library

CXX = g++
CXXFLAGS = -std=c++11 -fPIC -O3 -Wall
INCLUDES = -I/usr/include/opencv4
LIBS = -lopencv_core -lopencv_imgproc -lopencv_imgcodecs
TARGET = libimageprocessor.so
SOURCES = image_processor.cpp

all: $(TARGET)

$(TARGET): $(SOURCES)
	$(CXX) $(CXXFLAGS) $(INCLUDES) -shared -o $(TARGET) $(SOURCES) $(LIBS)

clean:
	rm -f $(TARGET) *.o *.a

check-deps:
	@echo "检查OpenCV安装..."
	@pkg-config --exists opencv4 && echo "✅ OpenCV4 已安装" || echo "❌ OpenCV4 未安装"

.PHONY: all clean check-deps 