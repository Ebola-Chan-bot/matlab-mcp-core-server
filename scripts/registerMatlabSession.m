function registerMatlabSession(sessionDir)
%REGISTERMATLABSESSION 注册当前 MATLAB 会话供 MCP 服务器发现
%
%   REGISTERMATLABSESSION() 在默认临时目录中注册当前会话。
%   MCP 服务器将自动发现并连接到此会话。
%
%   REGISTERMATLABSESSION(sessionDir) 在指定目录中注册会话。
%   如果目录不存在，将自动创建。
%
%   示例:
%       % 使用默认目录（推荐）
%       registerMatlabSession
%
%       % 使用自定义目录
%       registerMatlabSession('C:\MySessionDir')
%
%   注意:
%       - 需要确保 MATLAB 已启用 Embedded Connector
%       - MCP 服务器会在临时目录中搜索以 "matlab-mcp-core-server-" 开头的目录
%       - 会话目录名称以 "matlab-session-" 开头
%
%   See also: connector.securePort, mwapikey

    % 默认使用临时目录
    if nargin < 1 || isempty(sessionDir)
        tempDir = tempdir;
        % 使用用户名创建唯一的应用目录
        username = getenv('USERNAME');
        if isempty(username)
            username = getenv('USER');
        end
        if isempty(username)
            username = 'unknown';
        end
        appDir = fullfile(tempDir, ['matlab-mcp-core-server-manual']);
        sessionDir = fullfile(appDir, 'matlab-session-manual');
    end
    
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
    
    % 获取 API 密钥
    apiKey = getenv('MWAPIKEY');
    if isempty(apiKey)
        error('未设置 MWAPIKEY 环境变量。请在启动 MATLAB 前设置此变量。');
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
