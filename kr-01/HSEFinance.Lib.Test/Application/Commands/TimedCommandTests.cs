using HSEFinance.Lib.Application.Commands;
using Moq;
using System.Diagnostics;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Commands
{
    public class TimedCommandTests
    {
        [Fact]
        public void Execute_WhenCalled_MeasuresExecutionTime()
        {
            // Arrange
            var mockCommand = new Mock<ICommand>();
            Action<TimeSpan> onExecuted = elapsed => { /* обработка времени */ };
            var timedCommandMock = new Mock<Action<TimeSpan>>();
            var command = new TimedCommand(mockCommand.Object, timedCommandMock.Object);

            // Act
            command.Execute();

            // Assert
            mockCommand.Verify(c => c.Execute(), Times.Once);
            timedCommandMock.Verify(a => a(It.IsAny<TimeSpan>()), Times.Once);
        }

        [Fact]
        public void Execute_CorrectlyPassesElapsedTimeToAction()
        {
            // Arrange
            var mockCommand = new Mock<ICommand>();
            var elapsedDuration = TimeSpan.Zero;
            Action<TimeSpan> onExecuted = elapsed => elapsedDuration = elapsed;
            var command = new TimedCommand(mockCommand.Object, onExecuted);

            // Act
            var stopwatch = Stopwatch.StartNew();
            command.Execute();
            stopwatch.Stop();

            // Assert
            Assert.InRange(elapsedDuration.TotalMilliseconds, stopwatch.Elapsed.TotalMilliseconds - 1, stopwatch.Elapsed.TotalMilliseconds + 1);
        }
    }
}