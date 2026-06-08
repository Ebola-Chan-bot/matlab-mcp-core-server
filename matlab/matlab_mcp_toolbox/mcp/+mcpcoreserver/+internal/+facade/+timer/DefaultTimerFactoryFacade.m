classdef DefaultTimerFactoryFacade < mcpcoreserver.internal.facade.timer.TimerFactoryFacade
    %DefaultTimerFactoryFacade Creates DefaultTimerFacade instances

    % Copyright 2026 The MathWorks, Inc.

    methods
        function t = timer(~, period, executionMode, busyMode, timerFcn)
            t = mcpcoreserver.internal.facade.timer.DefaultTimerFacade( ...
                period, executionMode, busyMode, timerFcn);
        end
    end

end
