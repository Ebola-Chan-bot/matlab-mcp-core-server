function strategy = createDotNetAccessControl(options)
    %createDotNetAccessControl Factory that selects the correct strategy for the active .NET runtime

    % Copyright 2026 The MathWorks, Inc.

    arguments
        options.OSFacade mcpcoreserver.internal.facade.os.OSFacade = ...
            mcpcoreserver.internal.facade.os.DefaultOSFacade()
        options.DotNetFacade mcpcoreserver.internal.facade.dotnet.DotNetFacade = ...
            mcpcoreserver.internal.facade.dotnet.DefaultDotNetFacade()
    end

    if ~options.OSFacade.ispc()
        strategy = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.internal.dotnet.CoreDotNetAccessControl( ...
            DotNetFacade=options.DotNetFacade);
        return;
    end

    runtime = options.DotNetFacade.getDotNetRuntime();
    if runtime == "framework"
        strategy = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.internal.dotnet.FrameworkDotNetAccessControl();
    else
        strategy = mcpcoreserver.internal.fs.internal.permissionmanager.internal.windowsacl.internal.dotnet.CoreDotNetAccessControl( ...
            DotNetFacade=options.DotNetFacade);
    end
end
