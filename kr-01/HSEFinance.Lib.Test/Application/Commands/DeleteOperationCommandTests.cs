using HSEFinance.Lib.Application.Commands;
using HSEFinance.Lib.Domain.Repositories;
using Moq;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Commands
{
    public class DeleteOperationCommandTests
    {
        [Fact]
        public void Execute_WhenCalled_DeletesOperation()
        {
            // Arrange
            var operationRepositoryMock = new Mock<IOperationRepository>();
            var operationId = Guid.NewGuid();
            var command = new DeleteOperationCommand(operationRepositoryMock.Object, operationId);

            // Act
            command.Execute();

            // Assert
            operationRepositoryMock.Verify(r => r.DeleteOperation(operationId), Times.Once);
        }
    }
}