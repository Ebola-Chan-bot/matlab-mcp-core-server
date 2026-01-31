function registerMatlabSession
%REGISTERMATLABSESSION 注册当前 MATLAB 会话供 MCP 服务器发现
%
%   REGISTERMATLABSESSION() 在临时目录中注册当前会话。
%   MCP 服务器将自动发现并连接到此会话。
%
%   示例:
%       registerMatlabSession
%
%   注意:
%       - 需要确保 MATLAB 已启用 Embedded Connector
%       - MWAPIKEY 环境变量由 MATLAB 启动时自动设置
%
%   See also: connector.securePort

    % 使用固定的临时目录
    tempDir = tempdir;
    appDir = fullfile(tempDir, 'matlab-mcp-core-server-manual');
    sessionDir = fullfile(appDir, 'matlab-session-manual');
    
    % 确保目录存在
    if ~exist(sessionDir, 'dir')
        mkdir(sessionDir);
        fprintf('创建会话目录: %s\n', sessionDir);
    end
    
    % 获取 Embedded Connector 的安全端口
    try
        port = connector.securePort;
    catch ME
        error('无法获取 Embedded Connector 端口。请确保 MATLAB 已启用 Embedded Connector。\n错误: %s', ME.message);
    end
    
    if isempty(port) || port == 0
        error('Embedded Connector 未运行或未启用安全端口。');
    end
    
    % 获取 API 密钥（MATLAB 启动时自动设置）
    apiKey = getenv('MWAPIKEY');
    if isempty(apiKey)
        error('无法获取 MWAPIKEY。此环境变量应由 MATLAB 启动时自动设置。');
    end
    
    % 写入端口文件
    portFile = fullfile(sessionDir, 'connector.securePort');
    fid = fopen(portFile, 'w');
    if fid == -1
        error('无法创建端口文件: %s', portFile);
    end
    fprintf(fid, '%d', port);
    fclose(fid);
    fprintf('已写入端口文件: %s (端口: %d)\n', portFile, port);
    
    % 写入 API 密钥文件
    apiKeyFile = fullfile(sessionDir, 'apikey');
    fid = fopen(apiKeyFile, 'w');
    if fid == -1
        error('无法创建 API 密钥文件: %s', apiKeyFile);
    end
    fprintf(fid, '%s', apiKey);
    fclose(fid);
    fprintf('已写入 API 密钥文件: %s\n', apiKeyFile);
    
    % 显示成功消息
    fprintf('\n========================================\n');
    fprintf('MATLAB 会话已注册!\n');
    fprintf('========================================\n');
    fprintf('会话目录: %s\n', sessionDir);
    fprintf('端口: %d\n', port);
    fprintf('API 密钥: %s\n', apiKey);
    fprintf('\nMCP 服务器现在可以发现并连接到此会话。\n');
    fprintf('保持此 MATLAB 窗口打开以维持连接。\n');
    fprintf('========================================\n');
end
