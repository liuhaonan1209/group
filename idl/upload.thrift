namespace go car.upload

// ==================== 图片上传接口 ====================

// 图片上传请求
struct ImageUploadReq {
    1: string fileName          // 文件名
    2: binary fileData          // 文件数据
    3: i64 fileSize             // 文件大小
    4: string fileType          // 文件类型（扩展名）
    5: optional string uploadType // 上传类型（local/oss，默认oss）
}

// 图片上传响应
struct ImageUploadResp {
    1: string fileURL           // 文件访问URL
    2: string fileName          // 保存的文件名
    3: i64 fileSize             // 文件大小
    4: string message           // 返回消息
    5: bool success             // 是否成功
}

// ==================== 图片删除接口 ====================

// 图片删除请求
struct ImageDeleteReq {
    1: string fileURL           // 文件URL或文件路径
    2: optional string deleteType // 删除类型（local/oss，默认oss）
}

// 图片删除响应
struct ImageDeleteResp {
    1: string message           // 返回消息
    2: bool success             // 是否成功
}

// ==================== 批量图片上传接口 ====================

// 批量图片上传请求
struct BatchImageUploadReq {
    1: list<ImageUploadReq> images // 图片列表
}

// 批量上传结果
struct BatchUploadResult {
    1: string fileName          // 文件名
    2: string fileURL           // 文件URL
    3: bool success             // 是否成功
    4: string message           // 消息
}

// 批量图片上传响应
struct BatchImageUploadResp {
    1: list<BatchUploadResult> results // 上传结果列表
    2: i32 successCount         // 成功数量
    3: i32 failCount            // 失败数量
    4: string message           // 返回消息
    5: bool success             // 是否全部成功
}

// ==================== 图片信息查询接口 ====================

// 图片信息查询请求
struct ImageInfoReq {
    1: string fileURL           // 文件URL
}

// 图片信息响应
struct ImageInfoResp {
    1: string fileName          // 文件名
    2: i64 fileSize             // 文件大小
    3: string fileType          // 文件类型
    4: string uploadTime        // 上传时间
    5: string message           // 返回消息
    6: bool success             // 是否成功
}

// ==================== 服务定义 ====================

service UploadService {
    // 图片上传接口
    ImageUploadResp ImageUpload(1: ImageUploadReq req)
    
    // 图片删除接口
    ImageDeleteResp ImageDelete(1: ImageDeleteReq req)
    
    // 批量图片上传接口
    BatchImageUploadResp BatchImageUpload(1: BatchImageUploadReq req)
    
    // 图片信息查询接口
    ImageInfoResp ImageInfo(1: ImageInfoReq req)
}
