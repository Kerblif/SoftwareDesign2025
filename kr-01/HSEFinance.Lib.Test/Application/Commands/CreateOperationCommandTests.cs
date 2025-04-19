using HSEFinance.Lib.Application.Commands;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using Moq;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Commands
{
    public class CreateOperationCommandTests
    {
        [Fact]
        public void Execute_WhenCalled_InvokesRepositoryCreateOperation()
        {
            var operationRepositoryMock = new Mock<IOperationRepository>();
            var bankAccountId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();
            var command = new CreateOperationCommand(
                operationRepositoryMock.Object,
                ItemType.Income,
                bankAccountId,
                100.5m,
                DateTime.Now,
                categoryId,
                "Test Description"
            );

            command.Execute();

            operationRepositoryMock.Verify(r => r.CreateOperation(
                ItemType.Income,
                bankAccountId,
                100.5m,
                It.IsAny<DateTime>(),
                categoryId,
                "Test Description"
            ), Times.Once);
        }
    }
}