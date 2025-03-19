using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Repositories;
using Microsoft.EntityFrameworkCore;

namespace HSEFinance.Lib.Infrastructure.Data
{
    public class AccountRepository : IAccountRepository
    {
        private readonly HSEFinanceDbContext _dbContext;

        public AccountRepository(HSEFinanceDbContext dbContext)
        {
            _dbContext = dbContext;
        }

        public IEnumerable<BankAccount> GetAllBankAccounts()
        {
            return _dbContext.BankAccounts;
        }

        public BankAccount CreateBankAccount(string name)
        {
            var account = new BankAccount(name);
            _dbContext.BankAccounts.Add(account);
            _dbContext.SaveChanges();
            return account;
        }

        public BankAccount? GetBankAccount(Guid accountId)
        {
            return _dbContext.BankAccounts.Find(accountId);
        }

        public bool DeleteBankAccount(Guid accountId)
        {
            var account = GetBankAccount(accountId);
            if (account == null)
                return false;

            _dbContext.BankAccounts.Remove(account);
            _dbContext.SaveChanges();
            return true;
        }

        public void UpdateBankAccount(BankAccount account)
        {
            _dbContext.BankAccounts.Update(account);
            _dbContext.SaveChanges();
        }
    }
}