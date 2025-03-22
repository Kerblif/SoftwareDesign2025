using HSEFinance.Lib.Core;
using HSEFinance.Lib.Core.Interfaces;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Infrastructure.Data.Proxies
{
    public class AccountRepositoryProxy : IAccountRepository
    {
        private readonly IAccountRepository _repository;
        private readonly Dictionary<Guid, BankAccount> _cache = new();

        public AccountRepositoryProxy(IAccountRepository repository)
        {
            _repository = repository;

            var allAccounts = _repository.GetAllBankAccounts();
            foreach (var account in allAccounts)
            {
                _cache[account.Id] = account;
            }
        }

        public BankAccount CreateBankAccount(string? name)
        {
            var account = _repository.CreateBankAccount(name);

            _cache[account.Id] = account;

            return account;
        }

        public BankAccount? GetBankAccount(Guid accountId)
        {
            if (_cache.TryGetValue(accountId, out var account))
                return account;

            account = _repository.GetBankAccount(accountId);
            if (account != null)
            {
                _cache[accountId] = account;
            }

            return account;
        }

        public IEnumerable<BankAccount> GetAllBankAccounts()
        {
            return _cache.Values;
        }

        public bool DeleteBankAccount(Guid accountId)
        {
            var success = _repository.DeleteBankAccount(accountId);

            if (success)
            {
                _cache.Remove(accountId);
            }

            return success;
        }

        public void UpdateBankAccount(BankAccount account)
        {
            _repository.UpdateBankAccount(account);

            _cache[account.Id] = account;
        }

        public void UploadBankAccount(BankAccount account)
        {
            if (_cache.ContainsKey(account.Id))
            {
                throw new InvalidOperationException($"A bank account with ID {account.Id} already exists in the cache.");
            }
        
            _repository.UploadBankAccount(account);
            _cache[account.Id] = account;
        }

        public void RecalculateAccountBalance(Guid accountId)
        {
            if (!_cache.TryGetValue(accountId, out var cachedAccount))
            {
                cachedAccount = _repository.GetBankAccount(accountId) 
                                ?? throw new Exception($"Account with ID {accountId} not found.");
            }

            _repository.RecalculateAccountBalance(accountId);

            var recalculatedBalance = _repository.GetBankAccount(accountId)?.Balance
                                      ?? throw new Exception($"Failed to get recalculated balance for account {accountId}.");
            cachedAccount.Balance = recalculatedBalance;

            _cache[accountId] = cachedAccount;
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