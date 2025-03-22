using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;
using Moq;
using Xunit;

namespace HSEFinance.Lib.Test.Infrastructure.Data.Proxies
{
    public class OperationRepositoryProxyTests
    {
        [Fact]
        public void GetAllOperations_InitiallyLoadsAllOperationsFromBaseRepository()
        {
            // Arrange
            var mockRepository = new Mock<IOperationRepository>();
            var operations = new List<Operation>
            {
                new Operation(ItemType.Income, Guid.NewGuid(), 100m, DateTime.UtcNow, Guid.NewGuid()),
                new Operation(ItemType.Expense, Guid.NewGuid(), 50m, DateTime.UtcNow, Guid.NewGuid())
            };
            mockRepository.Setup(r => r.GetAllOperations()).Returns(operations);

            var proxy = new OperationRepositoryProxy(mockRepository.Object);

            // Act
            var allOperations = proxy.GetAllOperations();

            // Assert
            Assert.Equal(2, allOperations.Count());
            mockRepository.Verify(r => r.GetAllOperations(), Times.Once);
        }

        [Fact]
        public void GetOperation_UsesCacheForPreviouslyLoadedOperation()
        {
            // Arrange
            var mockRepository = new Mock<IOperationRepository>();
            var operation = new Operation(ItemType.Income, Guid.NewGuid(), 200m, DateTime.UtcNow, Guid.NewGuid());
            mockRepository.Setup(r => r.GetAllOperations()).Returns(new List<Operation> { operation });

            var proxy = new OperationRepositoryProxy(mockRepository.Object);

            // Act
            var retrievedOperation = proxy.GetOperation(operation.Id);

            // Assert
            Assert.NotNull(retrievedOperation);
            Assert.Equal(operation.Amount, retrievedOperation!.Amount);
            mockRepository.Verify(r => r.GetOperation(It.IsAny<Guid>()), Times.Never);
        }

        [Fact]
        public void CreateOperation_AddsToRepositoryAndCache()
        {
            // Arrange
            var mockRepository = new Mock<IOperationRepository>();
            var operation = new Operation(ItemType.Expense, Guid.NewGuid(), 300m, DateTime.UtcNow, Guid.NewGuid());
            mockRepository.Setup(r =>
                r.CreateOperation(operation.Type, operation.BankAccountId, operation.Amount, operation.Date, operation.CategoryId, null))
                .Returns(operation);

            var proxy = new OperationRepositoryProxy(mockRepository.Object);

            // Act
            var createdOperation = proxy.CreateOperation(operation.Type, operation.BankAccountId, operation.Amount, operation.Date, operation.CategoryId);

            // Assert
            Assert.NotNull(createdOperation);
            Assert.Equal(operation.Amount, createdOperation.Amount);
            mockRepository.Verify(r => r.CreateOperation(operation.Type, operation.BankAccountId, operation.Amount, operation.Date, operation.CategoryId, null), Times.Once);

            var cachedOperation = proxy.GetOperation(operation.Id);
            Assert.NotNull(cachedOperation);
            Assert.Equal(operation.Amount, cachedOperation!.Amount);
        }

        [Fact]
        public void DeleteOperation_RemovesFromCacheAndRepository()
        {
            // Arrange
            var mockRepository = new Mock<IOperationRepository>();
            var operation = new Operation(ItemType.Expense, Guid.NewGuid(), 50m, DateTime.UtcNow, Guid.NewGuid());
            mockRepository.Setup(r => r.GetAllOperations()).Returns(new List<Operation> { operation });
            mockRepository.Setup(r => r.DeleteOperation(operation.Id)).Returns(true);

            var proxy = new OperationRepositoryProxy(mockRepository.Object);

            // Act
            var deleteResult = proxy.DeleteOperation(operation.Id);

            // Assert
            Assert.True(deleteResult);
            Assert.Null(proxy.GetOperation(operation.Id));
            mockRepository.Verify(r => r.DeleteOperation(operation.Id), Times.Once);
        }
    }
}