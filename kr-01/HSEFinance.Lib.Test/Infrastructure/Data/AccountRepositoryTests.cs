using System.Linq;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Domain.Factories;
using HSEFinance.Lib.Infrastructure.Data;
using Xunit;

namespace HSEFinance.Lib.Test.Infrastructure.Data
{
    public class AccountRepositoryTests
    {
        [Fact]
        public void CreateBankAccount_ShouldAddAccountToDatabase()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var accountFactory = new BankAccountFactory();
            var repository = new AccountRepository(dbContext, accountFactory);

            // Act
            var account = repository.CreateBankAccount("Primary Account");

            // Assert
            Assert.NotNull(account);
            Assert.Equal("Primary Account", account.Name);
            Assert.Single(dbContext.BankAccounts);
        }

        [Fact]
        public void GetBankAccount_WithValidId_ReturnsAccount()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var accountFactory = new BankAccountFactory();
            var repository = new AccountRepository(dbContext, accountFactory);

            var createdAccount = repository.CreateBankAccount("Test Account");

            // Act
            var retrievedAccount = repository.GetBankAccount(createdAccount.Id);

            // Assert
            Assert.NotNull(retrievedAccount);
            Assert.Equal("Test Account", retrievedAccount?.Name);
        }

        [Fact]
        public void DeleteBankAccount_RemovesAccountFromDatabase()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var accountFactory = new BankAccountFactory();
            var repository = new AccountRepository(dbContext, accountFactory);

            var account = repository.CreateBankAccount("To Be Deleted");

            // Act
            var result = repository.DeleteBankAccount(account.Id);

            // Assert
            Assert.True(result);
            Assert.Empty(dbContext.BankAccounts);
        }

        [Fact]
        public void RecalculateAccountBalance_UpdatesBalanceCorrectly()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var accountFactory = new BankAccountFactory();
            var repository = new AccountRepository(dbContext, accountFactory);

            var account = repository.CreateBankAccount("Balance Test");
            dbContext.Operations.AddRange(
                new Operation(ItemType.Income, account.Id, 100m, DateTime.UtcNow, Guid.NewGuid()),
                new Operation(ItemType.Expense, account.Id, 50m, DateTime.UtcNow, Guid.NewGuid())
            );
            dbContext.SaveChanges();

            // Act
            repository.RecalculateAccountBalance(account.Id);
            var updatedAccount = repository.GetBankAccount(account.Id);

            // Assert
            Assert.NotNull(updatedAccount);
            Assert.Equal(50m, updatedAccount?.Balance);
        }
    }
}