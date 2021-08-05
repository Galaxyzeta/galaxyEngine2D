# 代码复习

## 启动流程

- 主线程是渲染线程。
- application.Start():
  - 执行用户自定义 init 函数
    - 注册System
    - 注册图像资源等
  - 启动线程池
  - 启动一个核心逻辑线程
  - 启动渲染无限循环
- 逻辑执行线程：
  - 创建对象并加入对象池
  - 执行 ECS 系统
  - 遍历对象池执行 step 操作
  - 输入状态清除
  - 执行对象销毁
  - 记忆对象当前 step 状态

## 新功能

- ID 生成器
- 判断某处是否有指定 name 的物体
- 丰富了碰撞检测方法

## TODO

- debug 模式启动有高概率导致并发访问map的问题