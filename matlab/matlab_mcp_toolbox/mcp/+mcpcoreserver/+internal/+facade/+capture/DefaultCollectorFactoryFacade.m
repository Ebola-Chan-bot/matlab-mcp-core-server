classdef DefaultCollectorFactoryFacade < mcpcoreserver.internal.facade.capture.CollectorFactoryFacade
    %DefaultCollectorFactoryFacade Creates DefaultCollectorFacade instances

    % Copyright 2026 The MathWorks, Inc.

    methods
        function collector = EventsCollector(~)
            collector = mcpcoreserver.internal.facade.capture.DefaultCollectorFacade();
        end
    end

end
