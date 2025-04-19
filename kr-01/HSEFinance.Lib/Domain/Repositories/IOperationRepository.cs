using HSEFinance.Lib.Core.Interfaces;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Domain.Repositories
{
    public interface IOperationRepository: IVisitable
    {
        IEnumerable<Operation> GetAllOperations();
        Operation CreateOperation(ItemType type,
                                  Guid bankAccountId,
                                  decimal amount,
                                  DateTime date,
                                  Guid categoryId,
                                  string? description = null);
        Operation? GetOperation(Guid operationId);
        bool DeleteOperation(Guid operationId);
        void UpdateOperation(Operation operation);
        void UploadOperation(Operation operation);
    }
}