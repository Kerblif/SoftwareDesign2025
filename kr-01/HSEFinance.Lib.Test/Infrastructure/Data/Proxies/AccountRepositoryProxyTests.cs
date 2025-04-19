using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;
using Moq;
using Xunit;

namespace HSEFinance.Lib.Test.Infrastructure.Data.Proxies
{
    public class AccountRepositoryProxyTests
    {
        [Fact]
        public void GetAllBankAccounts_InitiallyLoadsAllAccountsFromBaseRepository()
        {
            var mockRepository = new Mock<IAccountRepository>();
            var accounts = new List<BankAccount>
            {
                new BankAccount("Account 1"),
                new BankAccount("Account 2")
            };
            mockRepository.Setup(r => r.GetAllBankAccounts()).Returns(accounts);

            var proxy = new AccountRepositoryProxy(mockRepository.Object);

            var allAccounts = proxy.GetAllBankAccounts();

            Assert.Equal(2, allAccounts.Count());
            mockRepository.Verify(r => r.GetAllBankAccounts(), Times.Once);
        }

        [Fact]
        public void GetBankAccount_UsesCacheForPreviouslyLoadedAccount()
        {
            var mockRepository = new Mock<IAccountRepository>();
            var account = new BankAccount("Cached Account");
            mockRepository.Setup(r => r.GetAllBankAccounts()).Returns(new List<BankAccount> { });
            mockRepository.Setup(r => r.GetBankAccount(account.Id)).Throws(new Exception("Should not call repository!"));

            var proxy = new AccountRepositoryProxy(mockRepository.Object);

            proxy.UploadBankAccount(account);
            var retrievedAccount = proxy.GetBankAccount(account.Id);

            Assert.NotNull(retrievedAccount);
            Assert.Equal(account.Id, retrievedAccount.Id);
            mockRepository.Verify(r => r.GetBankAccount(It.IsAny<Guid>()), Times.Never);
        }

        [Fact]
        public void CreateBankAccount_AddsToRepositoryAndCache()
        {
            var mockRepository = new Mock<IAccountRepository>();
            var account = new BankAccount("New Account");
            mockRepository.Setup(r => r.CreateBankAccount("New Account")).Returns(account);

            var proxy = new AccountRepositoryProxy(mockRepository.Object);

            var createdAccount = proxy.CreateBankAccount("New Account");

            Assert.NotNull(createdAccount);
            Assert.Equal("New Account", createdAccount.Name);
            mockRepository.Verify(r => r.CreateBankAccount("New Account"), Times.Once);

            var cachedAccount = proxy.GetBankAccount(account.Id);
            Assert.Equal("New Account", cachedAccount!.Name);
        }

        [Fact]
        public void DeleteBankAccount_RemovesFromRepositoryAndCache()
        {
            var mockRepository = new Mock<IAccountRepository>();
            var accounts = new List<BankAccount>
            {
                new BankAccount("Account to Delete")
            };
            mockRepository.Setup(r => r.GetAllBankAccounts()).Returns(accounts);
            mockRepository.Setup(r => r.DeleteBankAccount(accounts[0].Id)).Returns(true);

            var proxy = new AccountRepositoryProxy(mockRepository.Object);

            var deleteResult = proxy.DeleteBankAccount(accounts[0].Id);

            Assert.True(deleteResult);
            Assert.Null(proxy.GetBankAccount(accounts[0].Id)); // Cache should be updated
            mockRepository.Verify(r => r.DeleteBankAccount(accounts[0].Id), Times.Once);
        }
    }
}