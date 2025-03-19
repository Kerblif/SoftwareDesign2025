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

        public BankAccount CreateBankAccount(string name)
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
    }
}