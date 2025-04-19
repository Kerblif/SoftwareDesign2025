using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Application.Commands
{
    public class CreateOperationCommand : ICommand
    {
        private readonly IOperationRepository _operationRepository;
        private readonly ItemType _type;
        private readonly Guid _bankAccountId;
        private readonly decimal _amount;
        private readonly DateTime _date;
        private readonly Guid _categoryId;
        private readonly string? _description;

        public CreateOperationCommand(IOperationRepository operationRepository, ItemType type, Guid bankAccountId, decimal amount, DateTime date, Guid categoryId, string? description = null)
        {
            _operationRepository = operationRepository;
            _type = type;
            _bankAccountId = bankAccountId;
            _amount = amount;
            _date = date;
            _categoryId = categoryId;
            _description = description;
        }

        public void Execute()
        {
            _operationRepository.CreateOperation(_type, _bankAccountId, _amount, _date, _categoryId, _description);
        }
    }
}