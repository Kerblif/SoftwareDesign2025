using HSEFinance.Lib.Application.Facades;

namespace HSEFinance.Lib.Application.Commands
{
    public class DeleteOperationCommand : ICommand
    {
        private readonly OperationFacade _operationFacade;
        private readonly Guid _operationId;

        public DeleteOperationCommand(OperationFacade operationFacade, Guid operationId)
        {
            _operationFacade = operationFacade;
            _operationId = operationId;
        }

        public void Execute()
        {
            _operationFacade.DeleteOperation(_operationId);
        }
    }
}