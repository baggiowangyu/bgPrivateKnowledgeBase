# OpenCV编译的套路

## Windows环境编译

- 首先到GitHub上克隆几个项目：
  - opencv
  - opencv_contrib
  - opencv_extra
  由于4.0版本开始opencv需要用C++11来编译，于是我选择了切换到3.4.5版本

- 使用CMake，生成VS2013版本工程配置
  - 实际上这里使用2013以前的版本是会悲剧的，因为OpenCV的thirdparty中的quirc项目使用了C99，VS2013以前的版本会遇到编译不过的情况
  - 这里我也是摸了许久才摸透...

## 关于人脸库的建设

- 实际上我们要为每一个人脸进行学习入库
  - 训练的结果保存在xml文件中，以数值id标记
  - 实际的人脸数据库应该与这个数值id实现关联，人脸数据库字段[人员ID | 人员姓名 | 性别 | 出生日期 | 民族 | ]
