using HSEFinance.Lib.Domain.Repositories;
namespace HSEFinance.Lib.Application.Commands
{
    public class DeleteOperationCommand : ICommand
    {
        private readonly IOperationRepository _operationRepository;
        private readonly Guid _operationId;

        public DeleteOperationCommand(IOperationRepository operationRepository, Guid operationId)
        {
            _operationRepository = operationRepository;
            _operationId = operationId;
        }

        public void Execute()
        {
            _operationRepository.DeleteOperation(_operationId);
        }
    }
}