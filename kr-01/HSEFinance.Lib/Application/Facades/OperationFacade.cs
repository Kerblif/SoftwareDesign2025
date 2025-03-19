using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Application.Facades
{
    public class OperationFacade
    {
        private readonly IAccountRepository _accountRepository;
        private readonly IOperationFactory _operationFactory;

        public OperationFacade(IAccountRepository accountRepository, IOperationFactory operationFactory)
        {
            _accountRepository = accountRepository ?? throw new ArgumentNullException(nameof(accountRepository));
            _operationFactory = operationFactory ?? throw new ArgumentNullException(nameof(operationFactory));
        }

        public Operation CreateOperation(
            ItemType type,
            Guid bankAccountId,
            decimal amount,
            DateTime date,
            Guid categoryId,
            string? description = null)
        {
            var operation = _operationFactory.Create(type, bankAccountId, amount, date, categoryId, description);

            var account = _accountRepository.GetBankAccount(bankAccountId);
            if (account != null)
            {
                account.UpdateBalance(amount, type);
                _accountRepository.UpdateBankAccount(account);
            }

            return operation;
        }

        public Operation? GetOperation(Guid operationId)
        {
            throw new NotImplementedException("You should implement storage and retrieval logic.");
        }

        public bool DeleteOperation(Guid operationId)
        {
            throw new NotImplementedException("You should implement delete logic.");
        }
    }
}