using HSEFinance.Lib.Core;
using HSEFinance.Lib.Core.Interfaces;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
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

        public BankAccount CreateBankAccount(string? name)
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

        public void UploadBankAccount(BankAccount account)
        {
            if (_dbContext.BankAccounts.Find(account.Id) != null)
            {
                throw new InvalidOperationException($"A bank account with ID {account.Id} already exists.");
            }

            _dbContext.BankAccounts.Add(account);
            _dbContext.SaveChanges();
        }

        public void RecalculateAccountBalance(Guid accountId)
        {
            var bankAccount = GetBankAccount(accountId);
            if (bankAccount == null)
            {
                throw new Exception($"Bank account with ID {accountId} not found.");
            }
        
            bankAccount.Balance = _dbContext.Operations
                .Where(o => o.BankAccountId == accountId).Sum(o => o.Type == ItemType.Income ? o.Amount : -o.Amount);
            
            _dbContext.BankAccounts.Update(bankAccount);
            _dbContext.SaveChanges();
        }

        public void Accept(IVisitor visitor)
        {
            foreach (var account in GetAllBankAccounts())
            {
                account.Accept(visitor);
            }
        }
    }
}