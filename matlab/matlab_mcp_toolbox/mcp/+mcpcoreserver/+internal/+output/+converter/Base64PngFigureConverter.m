classdef Base64PngFigureConverter < mcpcoreserver.internal.output.converter.Converter
    %Base64PngFigureConverter Converts figure handles to base64-encoded PNG

    % Copyright 2026 The MathWorks, Inc.

    properties (Constant)
        DefaultResolution(1, 1) double = 150
    end

    properties (GetAccess = private, SetAccess = immutable)
        AppDataLocator(1, 1) mcpcoreserver.internal.appdata.AppDataLocator = mcpcoreserver.internal.appdata.DefaultAppDataLocator()
        FSAdaptor(1, 1) mcpcoreserver.internal.fs.FSAdaptor = mcpcoreserver.internal.fs.DefaultFSAdaptor()
        ExportFacade(1, 1) mcpcoreserver.internal.facade.export.ExportFacade = mcpcoreserver.internal.facade.export.DefaultExportFacade()
        FSFacade(1, 1) mcpcoreserver.internal.facade.fs.FSFacade = mcpcoreserver.internal.facade.fs.DefaultFSFacade()
        Resolution(1, 1) double = mcpcoreserver.internal.output.converter.Base64PngFigureConverter.DefaultResolution
    end

    methods
        function obj = Base64PngFigureConverter(options)
            arguments
                options.?mcpcoreserver.internal.output.converter.Base64PngFigureConverter
            end

            for prop = string(fieldnames(options).')
                obj.(prop) = options.(prop);
            end
        end

        function out = convert(obj, evt)
            figHandle = evt.payload;

            figuresFolder = fullfile(obj.AppDataLocator.getAppDataFolder(), "v1", "figures");
            obj.FSAdaptor.ensureSecureFolder(figuresFolder);

            tmpPath = obj.ExportFacade.tempname(figuresFolder) + ".png";
            cleanup = onCleanup(@() obj.deleteTempFile(tmpPath));

            obj.ExportFacade.exportgraphics(figHandle, tmpPath, 'Resolution', obj.Resolution);

            bytes = obj.readFileBytes(tmpPath);
            b64 = obj.ExportFacade.base64encode(bytes);

            out = struct( ...
                'outputType', 'figure', ...
                'mimeType', 'image/png', ...
                'data', char(b64) ...
            );
        end
    end

    methods (Access = private)
        function bytes = readFileBytes(obj, path)
            fid = obj.FSFacade.fopen(path, 'r');
            if fid == -1
                throw(mcpcoreserver.internal.error.Errors.FailedToReadFile(path));
            end
            cleanup = onCleanup(@() obj.FSFacade.fclose(fid));
            bytes = obj.FSFacade.fread(fid, Inf, '*uint8');
        end

        function deleteTempFile(obj, path)
            try
                obj.FSFacade.delete(path);
            catch
            end
        end
    end

end
