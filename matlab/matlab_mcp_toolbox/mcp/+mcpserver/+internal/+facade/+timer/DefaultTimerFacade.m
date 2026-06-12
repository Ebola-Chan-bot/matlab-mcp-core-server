classdef DefaultTimerFacade < mcpserver.internal.facade.timer.TimerFacade
    %DefaultTimerFacade Wraps a single MATLAB timer object

    % Copyright 2026 The MathWorks, Inc.

    properties (Access = private)
        Timer
    end

    methods
        function obj = DefaultTimerFacade(period, executionMode, busyMode, timerFcn)
            obj.Timer = timer('Period', period, 'ExecutionMode', executionMode, ...
                'BusyMode', busyMode, 'TimerFcn', timerFcn);
        end

        function Start(obj)
            start(obj.Timer);
        end

        function Stop(obj)
            stop(obj.Timer);
        end

        function delete(obj)
            if ~isempty(obj.Timer) && isvalid(obj.Timer)
                stop(obj.Timer);
                delete(obj.Timer);
            end
        end
    end

end
