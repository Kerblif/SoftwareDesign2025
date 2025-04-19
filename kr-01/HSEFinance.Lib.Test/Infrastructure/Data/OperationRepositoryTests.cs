using System.Linq;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Domain.Factories;
using HSEFinance.Lib.Infrastructure.Data;
using Xunit;

namespace HSEFinance.Lib.Test.Infrastructure.Data
{
    public class OperationRepositoryTests
    {
        static HSEFinanceDbContext dbContext = HSEFinanceDbContextFactory.Create();
        static BankAccountFactory accountFactory = new BankAccountFactory();
        static CategoryFactory categoryFactory = new CategoryFactory();
        static OperationFactory operationFactory = new OperationFactory();
        static AccountRepository accountRepository = new AccountRepository(dbContext, accountFactory);
        static ICategoryRepository categoryRepository = new CategoryRepository(dbContext, categoryFactory);
        static IOperationRepository repository = new OperationRepository(dbContext, operationFactory, accountRepository, categoryRepository);
        
        [Fact]
        public void CreateOperation_ShouldAddOperationAndUpdateAccountBalance()
        {
            // Arrange
            var account = accountRepository.CreateBankAccount("Test Account");
            var category = categoryRepository.CreateCategory(ItemType.Income, "Salary");

            // Act
            var operation = repository.CreateOperation(ItemType.Income, account.Id, 500, DateTime.UtcNow, category.Id);

            // Assert
            Assert.NotNull(operation);
            Assert.Equal(500, accountRepository.GetBankAccount(account.Id)?.Balance);
            Assert.Single(dbContext.Operations);
        }

        [Fact]
        public void DeleteOperation_ShouldRemoveOperationAndUpdateAccountBalance()
        {
            // Arrange
            var account = accountRepository.CreateBankAccount("Test Account");
            var category = categoryRepository.CreateCategory(ItemType.Expense, "Food");

            var operation = repository.CreateOperation(ItemType.Expense, account.Id, 100, DateTime.UtcNow, category.Id);

            // Act
            repository.DeleteOperation(operation.Id);

            // Assert
            Assert.Equal(0, accountRepository.GetBankAccount(account.Id)?.Balance);
            Assert.Empty(dbContext.Operations);
        }

        [Fact]
        public void CreateOperation_WithInvalidAmount_ThrowsArgumentException()
        {
            // Arrange
            var bankAccountId = Guid.NewGuid();
            var categoryId = Guid.NewGuid();

            Assert.Throws<ArgumentException>(() =>
            {
                // Operation with negative amount
                var operation = new Operation(
                    ItemType.Expense,
                    bankAccountId,
                    -100.0m,
                    DateTime.Now,
                    categoryId,
                    null
                );
            });
        }

        [Fact]
        public void CreateOperation_WithNullCategory_ThrowsException()
        {
            // Arrange
            var bankAccountId = Guid.NewGuid();

            Assert.Throws<ArgumentException>(() =>
            {
                var operation = new Operation(
                    ItemType.Income,
                    bankAccountId,
                    150.0m,
                    DateTime.Now,
                    Guid.Empty, // Invalid empty CategoryId
                    null
                );
            });
        }
    }
}