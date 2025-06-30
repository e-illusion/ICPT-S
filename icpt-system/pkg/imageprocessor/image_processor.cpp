// image_processor.cpp - 高性能C++图像处理模块
// 使用OpenCV实现快速图像压缩和缩略图生成

#include <opencv2/opencv.hpp>
#include <opencv2/imgproc.hpp>
#include <opencv2/imgcodecs.hpp>
#include <string>
#include <vector>
#include <memory>
#include <cstring>

extern "C" {

// 图像信息结构体
struct ImageInfo {
    int width;
    int height;
    int channels;
    size_t data_size;
    char format[16];
};

// 压缩配置结构体
struct CompressConfig {
    int quality;        // JPEG质量 (1-100)
    int max_width;      // 最大宽度
    int max_height;     // 最大高度
    bool enable_resize; // 是否启用尺寸调整
};

// 错误码定义
enum ErrorCode {
    SUCCESS = 0,
    ERROR_FILE_NOT_FOUND = -1,
    ERROR_INVALID_IMAGE = -2,
    ERROR_SAVE_FAILED = -3,
    ERROR_MEMORY_ALLOCATION = -4,
    ERROR_INVALID_PARAMS = -5
};

// 高性能图像压缩函数
int compress_image(const char* input_path, const char* output_path, const CompressConfig* config) {
    try {
        if (!input_path || !output_path || !config) {
            return ERROR_INVALID_PARAMS;
        }

        // 读取图像
        cv::Mat image = cv::imread(input_path, cv::IMREAD_COLOR);
        if (image.empty()) {
            return ERROR_FILE_NOT_FOUND;
        }

        // 尺寸调整（如果启用）
        if (config->enable_resize && 
            (image.cols > config->max_width || image.rows > config->max_height)) {
            
            double scale_w = static_cast<double>(config->max_width) / image.cols;
            double scale_h = static_cast<double>(config->max_height) / image.rows;
            double scale = std::min(scale_w, scale_h);
            
            int new_width = static_cast<int>(image.cols * scale);
            int new_height = static_cast<int>(image.rows * scale);
            
            cv::Mat resized_image;
            cv::resize(image, resized_image, cv::Size(new_width, new_height), 
                      0, 0, cv::INTER_LANCZOS4);
            image = resized_image;
        }

        // 设置压缩参数
        std::vector<int> compression_params;
        compression_params.push_back(cv::IMWRITE_JPEG_QUALITY);
        compression_params.push_back(config->quality);
        compression_params.push_back(cv::IMWRITE_JPEG_OPTIMIZE);
        compression_params.push_back(1); // 启用优化

        // 保存压缩后的图像
        bool success = cv::imwrite(output_path, image, compression_params);
        if (!success) {
            return ERROR_SAVE_FAILED;
        }

        return SUCCESS;
    } catch (const std::exception& e) {
        return ERROR_INVALID_IMAGE;
    }
}

// 高性能缩略图生成函数
int generate_thumbnail(const char* input_path, const char* output_path, int thumb_width) {
    try {
        if (!input_path || !output_path || thumb_width <= 0) {
            return ERROR_INVALID_PARAMS;
        }

        // 读取图像
        cv::Mat image = cv::imread(input_path, cv::IMREAD_COLOR);
        if (image.empty()) {
            return ERROR_FILE_NOT_FOUND;
        }

        // 计算缩略图尺寸，保持宽高比
        double aspect_ratio = static_cast<double>(image.rows) / image.cols;
        int thumb_height = static_cast<int>(thumb_width * aspect_ratio);

        // 生成缩略图（使用高质量插值）
        cv::Mat thumbnail;
        cv::resize(image, thumbnail, cv::Size(thumb_width, thumb_height), 
                  0, 0, cv::INTER_LANCZOS4);

        // 应用锐化滤波器以提升质量
        cv::Mat sharpened;
        cv::Mat kernel = (cv::Mat_<float>(3,3) << 
                         0, -1, 0,
                        -1,  5, -1,
                         0, -1, 0);
        cv::filter2D(thumbnail, sharpened, thumbnail.depth(), kernel);

        // 设置高质量JPEG压缩参数
        std::vector<int> compression_params;
        compression_params.push_back(cv::IMWRITE_JPEG_QUALITY);
        compression_params.push_back(95); // 高质量
        compression_params.push_back(cv::IMWRITE_JPEG_OPTIMIZE);
        compression_params.push_back(1);

        // 保存缩略图
        bool success = cv::imwrite(output_path, sharpened, compression_params);
        if (!success) {
            return ERROR_SAVE_FAILED;
        }

        return SUCCESS;
    } catch (const std::exception& e) {
        return ERROR_INVALID_IMAGE;
    }
}

// 获取图像信息函数
int get_image_info(const char* input_path, ImageInfo* info) {
    try {
        if (!input_path || !info) {
            return ERROR_INVALID_PARAMS;
        }

        // 读取图像信息（不加载完整数据）
        cv::Mat image = cv::imread(input_path, cv::IMREAD_UNCHANGED);
        if (image.empty()) {
            return ERROR_FILE_NOT_FOUND;
        }

        // 填充图像信息
        info->width = image.cols;
        info->height = image.rows;
        info->channels = image.channels();
        info->data_size = image.total() * image.elemSize();

        // 确定图像格式
        std::string ext = std::string(input_path);
        size_t dot_pos = ext.find_last_of('.');
        if (dot_pos != std::string::npos) {
            ext = ext.substr(dot_pos + 1);
            std::transform(ext.begin(), ext.end(), ext.begin(), ::tolower);
            strncpy(info->format, ext.c_str(), sizeof(info->format) - 1);
            info->format[sizeof(info->format) - 1] = '\0';
        } else {
            strcpy(info->format, "unknown");
        }

        return SUCCESS;
    } catch (const std::exception& e) {
        return ERROR_INVALID_IMAGE;
    }
}

// 批量图像处理函数
int batch_process_images(const char** input_paths, const char** output_paths, 
                        int count, const CompressConfig* config) {
    try {
        if (!input_paths || !output_paths || count <= 0 || !config) {
            return ERROR_INVALID_PARAMS;
        }

        int success_count = 0;
        for (int i = 0; i < count; i++) {
            if (compress_image(input_paths[i], output_paths[i], config) == SUCCESS) {
                success_count++;
            }
        }

        return success_count;
    } catch (const std::exception& e) {
        return ERROR_INVALID_IMAGE;
    }
}

// 内存中图像处理函数（高级功能）
int process_image_memory(const unsigned char* input_data, size_t input_size,
                        unsigned char** output_data, size_t* output_size,
                        const CompressConfig* config) {
    try {
        if (!input_data || input_size == 0 || !output_data || !output_size || !config) {
            return ERROR_INVALID_PARAMS;
        }

        // 从内存中解码图像
        std::vector<unsigned char> input_vector(input_data, input_data + input_size);
        cv::Mat image = cv::imdecode(input_vector, cv::IMREAD_COLOR);
        if (image.empty()) {
            return ERROR_INVALID_IMAGE;
        }

        // 应用压缩配置
        if (config->enable_resize && 
            (image.cols > config->max_width || image.rows > config->max_height)) {
            
            double scale_w = static_cast<double>(config->max_width) / image.cols;
            double scale_h = static_cast<double>(config->max_height) / image.rows;
            double scale = std::min(scale_w, scale_h);
            
            int new_width = static_cast<int>(image.cols * scale);
            int new_height = static_cast<int>(image.rows * scale);
            
            cv::Mat resized_image;
            cv::resize(image, resized_image, cv::Size(new_width, new_height), 
                      0, 0, cv::INTER_LANCZOS4);
            image = resized_image;
        }

        // 编码为JPEG
        std::vector<unsigned char> output_vector;
        std::vector<int> compression_params;
        compression_params.push_back(cv::IMWRITE_JPEG_QUALITY);
        compression_params.push_back(config->quality);
        compression_params.push_back(cv::IMWRITE_JPEG_OPTIMIZE);
        compression_params.push_back(1);

        bool success = cv::imencode(".jpg", image, output_vector, compression_params);
        if (!success) {
            return ERROR_SAVE_FAILED;
        }

        // 分配输出内存
        *output_size = output_vector.size();
        *output_data = static_cast<unsigned char*>(malloc(*output_size));
        if (!*output_data) {
            return ERROR_MEMORY_ALLOCATION;
        }

        memcpy(*output_data, output_vector.data(), *output_size);
        return SUCCESS;
    } catch (const std::exception& e) {
        return ERROR_INVALID_IMAGE;
    }
}

// 释放内存函数
void free_image_data(unsigned char* data) {
    if (data) {
        free(data);
    }
}

// 获取版本信息
const char* get_version() {
    return "ICPT C++ Image Processor v1.0.0";
}

// 获取OpenCV版本
const char* get_opencv_version() {
    static std::string version = cv::getBuildInformation();
    return version.c_str();
}

} // extern "C" 