classdef ExceptionReport
    %ExceptionReport Utilities for formatting MException reports

    % Copyright 2026 The MathWorks, Inc.

    methods (Static)
        function ex = stripInternalFrames(userError, callerName)
            %stripInternalFrames Build an MException with internal stack frames removed
            %   Creates a new MException whose message is the original error
            %   report with stack frames from callerName stripped out. If the
            %   caller frame is the only frame, the original message is used.
            arguments
                userError(1, 1) MException
                callerName(1, 1) string
            end

            report = getReport(userError, 'extended', 'hyperlinks', 'off');
            lines = splitlines(report);
            idx = find(startsWith(lines, "Error") & contains(lines, callerName), 1);
            if ~isempty(idx) && idx > 1
                report = strjoin(lines(1:idx-1), newline);
            else
                report = userError.message;
            end
            ex = MException(userError.identifier, '%s', report);
            for i = 1:numel(userError.cause)
                ex = ex.addCause(userError.cause{i});
            end
        end
    end

end
