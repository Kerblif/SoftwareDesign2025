using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Application.Facades
{
    public class BankAccountFacade
    {
        private readonly IAccountRepository _accountRepository;
        private readonly IBankAccountFactory _bankAccountFactory;

        public BankAccountFacade(IAccountRepository accountRepository, IBankAccountFactory bankAccountFactory)
        {
            _accountRepository = accountRepository ?? throw new ArgumentNullException(nameof(accountRepository));
            _bankAccountFactory = bankAccountFactory ?? throw new ArgumentNullException(nameof(bankAccountFactory));
        }

        public BankAccount CreateBankAccount(string name)
        {
            var account = _bankAccountFactory.Create(name);

            _accountRepository.CreateBankAccount(name);

            return account;
        }

        public BankAccount? GetBankAccount(Guid accountId)
        {
            return _accountRepository.GetBankAccount(accountId);
        }

        public bool DeleteBankAccount(Guid accountId)
        {
            return _accountRepository.DeleteBankAccount(accountId);
        }
    }
}