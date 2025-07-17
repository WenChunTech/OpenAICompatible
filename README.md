# OpenAI API Compatible

## 项目简介

本项目是一个与 OpenAI API 兼容的代理服务。它旨在提供一个与 OpenAI 官方 API 格式一致的接口，从而方便地将现有应用或生态工具无缝对接到各类大语言模型服务，而无需修改大量代码。

## 项目结构

项目代码主要位于 `src` 目录下，并遵循模块化的设计原则，各个模块职责清晰：

```
src/
├── config/      # 负责加载和管理项目配置
├── constant/    # 定义项目中使用的常量
├── converter/   # 负责在不同数据模型之间进行转换（例如，将特定模型的响应转换为 OpenAI 格式）
├── error/       # 定义和处理自定义错误
├── handler/     # 存放 HTTP 请求处理器，是 API 的业务逻辑核心
├── model/       # 定义项目中使用的数据结构，如 API 请求体和响应体
├── parser/      # 负责解析数据流或请求
├── request/     # 封装了向上游服务发起 HTTP 请求的逻辑
└── sse/         # 实现 Server-Sent Events (SSE)，用于支持流式 API 响应
```

## API 端点

服务启动后，会暴露以下兼容 OpenAI 格式的 API 端点：

*   `GET /v1/models`
    *   **功能**: 获取当前代理服务支持的模型列表。
    *   **处理程序**: `handler.ChatProxyModelHandler`
    *   **描述**: 返回一个包含多个模型信息的 JSON 数组，格式与 OpenAI 的 `v1/models` 接口一致。

*   `POST /v1/chat/completions`
    *   **功能**: 发起对话请求。
    *   **处理程序**: `handler.ChatProxyChatHandler`
    *   **描述**: 接收 OpenAI 格式的聊天请求，并代理到后端的语言模型服务。支持流式（`stream: true`）和非流式两种模式。

## 如何运行

您可以按照以下步骤在本地启动此服务：

1.  **克隆项目**
    ```bash
    git clone https://github.com/WenChunTech/OpenapiCompatible.git
    ```

2.  **进入项目目录**
    ```bash
    cd OpenapiCompatible
    ```

3.  **运行服务**
    ```bash
    go run main.go
    ```

4.  服务启动后，您将在控制台看到以下日志，并可以通过 `http://localhost:8080` 访问服务。
    ```
    Server starting on port 8080...
    ```

## 贡献

如果您有任何改进意见或想要贡献代码，请随时提交 Pull Request 或创建 Issue。提交规范请参考 [Commit Rule](./COMMIT_RULE.md)。

## 许可证

本项目采用 MIT 许可证。请查看 `LICENSE` 文件了解更多信息。