classdef registerMatlabSession < handle
	%REGISTERMATLABSESSION 注册当前 MATLAB 会话供 MCP 服务器发现（单例模式）
	%
	%   session = registerMatlabSession() 获取或创建会话注册对象。
	%   如果已存在注册对象，则返回现有实例；否则创建新实例。
	%   MCP 服务器将自动发现并连接到此会话。
	%
	%   此类使用单例模式，确保每个 MATLAB 会话只有一个注册实例。
	%   当对象被删除时（delete(registerMatlabSession) 或 MATLAB 退出），会话文件将自动清理。
	%
	%   示例:
	%       session = registerMatlabSession;  % 首次调用创建注册
	%       session2 = registerMatlabSession; % 返回同一实例
	%       % ... 使用 MCP 服务器 ...
	%       delete(registerMatlabSession)     % 取消注册并清理会话文件
	
	properties (SetAccess = immutable)
		SessionDir string   % 会话目录路径
		Port double         % Embedded Connector 端口
		APIKey string       % API 密钥
	end
	
	methods
		function obj = registerMatlabSession()
			%REGISTERMATLABSESSION 构造函数，创建会话文件
			SI=SingleInstance;
			if isempty(SI)
				% 使用固定的临时目录
				tempDir = tempdir;
				appDir = fullfile(tempDir, 'matlab-mcp-core-server-manual');
				obj.SessionDir = fullfile(appDir, 'matlab-session-manual');
				
				% 确保目录存在
				if ~exist(obj.SessionDir, 'dir')
					mkdir(obj.SessionDir);
					fprintf('创建会话目录: %s\n', obj.SessionDir);
				end
				
				% 获取 Embedded Connector 的安全端口
				try
					obj.Port = connector.securePort;
				catch ME
					error('registerMatlabSession:ConnectorError', ...
						'无法获取 Embedded Connector 端口。请确保 MATLAB 已启用 Embedded Connector。\n错误: %s', ME.message);
				end
				
				if isempty(obj.Port) || obj.Port == 0
					error('registerMatlabSession:NoSecurePort', ...
						'Embedded Connector 未运行或未启用安全端口。');
				end
				
				% 获取 API 密钥（MATLAB 启动时自动设置）
				obj.APIKey = getenv('MWAPIKEY');
				if isempty(obj.APIKey)
					error('registerMatlabSession:NoAPIKey', ...
						'无法获取 MWAPIKEY。此环境变量应由 MATLAB 启动时自动设置。');
				end
				
				% 写入端口文件
				portFile = fullfile(obj.SessionDir, 'connector.securePort');
				fid = fopen(portFile, 'w');
				if fid == -1
					error('registerMatlabSession:FileError', ...
						'无法创建端口文件: %s', portFile);
				end
				fprintf(fid, '%d', obj.Port);
				fclose(fid);
				fprintf('已写入端口文件: %s (端口: %d)\n', portFile, obj.Port);
				
				% 写入 API 密钥文件
				apiKeyFile = fullfile(obj.SessionDir, 'apikey');
				fid = fopen(apiKeyFile, 'w');
				if fid == -1
					error('registerMatlabSession:FileError', ...
						'无法创建 API 密钥文件: %s', apiKeyFile);
				end
				fprintf(fid, '%s', obj.APIKey);
				fclose(fid);
				fprintf('已写入 API 密钥文件: %s\n', apiKeyFile);
				
				% 显示成功消息
				fprintf('\n========================================\n');
				fprintf('MATLAB 会话已注册!\n');
				fprintf('========================================\n');
				fprintf('会话目录: %s\n', obj.SessionDir);
				fprintf('端口: %d\n', obj.Port);
				fprintf('API 密钥: %s\n', obj.APIKey);
				fprintf('\nMCP 服务器现在可以发现并连接到此会话。\n');
				fprintf('保持此对象以维持注册，clear 对象将取消注册。\n');
				fprintf('========================================\n');
				SingleInstance(obj);
			else
				obj=SI;
			end
		end
		
		function delete(obj)
			%DELETE 析构函数，删除会话文件
			
			if ~isempty(obj.SessionDir) && exist(obj.SessionDir, 'dir')
				try
					rmdir(obj.SessionDir, 's');
					fprintf('已清理会话目录: %s\n', obj.SessionDir);
				catch ME
					warning('registerMatlabSession:CleanupError', ...
						'清理会话目录失败: %s', ME.message);
				end
			end
		end
	end
end
function SI=SingleInstance(obj)
persistent pSI
if nargin
	pSI=obj;
end
SI=pSI;
end