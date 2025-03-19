using System.Diagnostics;

namespace HSEFinance.Lib.Application.Commands
{
    public class TimedCommand : ICommand
    {
        private readonly ICommand _command;
        private readonly Action<TimeSpan> _onExecuted;

        public TimedCommand(ICommand command, Action<TimeSpan> onExecuted)
        {
            _command = command;
            _onExecuted = onExecuted;
        }

        public void Execute()
        {
            var stopwatch = Stopwatch.StartNew();
            _command.Execute();
            stopwatch.Stop();

            _onExecuted(stopwatch.Elapsed);
        }
    }
}