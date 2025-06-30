// image_processor.h - C++图像处理库头文件
// 定义供Go语言调用的C接口

#ifndef IMAGE_PROCESSOR_H
#define IMAGE_PROCESSOR_H

#include <stddef.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// 图像信息结构体
typedef struct {
    int width;
    int height;
    int channels;
    size_t data_size;
    char format[16];
} ImageInfo;

// 压缩配置结构体
typedef struct {
    int quality;        // JPEG质量 (1-100)
    int max_width;      // 最大宽度
    int max_height;     // 最大高度
    bool enable_resize; // 是否启用尺寸调整
} CompressConfig;

// 错误码定义
typedef enum {
    SUCCESS = 0,
    ERROR_FILE_NOT_FOUND = -1,
    ERROR_INVALID_IMAGE = -2,
    ERROR_SAVE_FAILED = -3,
    ERROR_MEMORY_ALLOCATION = -4,
    ERROR_INVALID_PARAMS = -5
} ErrorCode;

// 函数声明
int compress_image(const char* input_path, const char* output_path, const CompressConfig* config);
int generate_thumbnail(const char* input_path, const char* output_path, int thumb_width);
int get_image_info(const char* input_path, ImageInfo* info);
int batch_process_images(const char** input_paths, const char** output_paths, int count, const CompressConfig* config);
int process_image_memory(const unsigned char* input_data, size_t input_size, unsigned char** output_data, size_t* output_size, const CompressConfig* config);
void free_image_data(unsigned char* data);
const char* get_version();
const char* get_opencv_version();

#ifdef __cplusplus
}
#endif

#endif // IMAGE_PROCESSOR_H 