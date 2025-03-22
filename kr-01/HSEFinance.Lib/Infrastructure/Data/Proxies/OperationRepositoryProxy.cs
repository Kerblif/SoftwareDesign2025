using HSEFinance.Lib.Core;
using HSEFinance.Lib.Core.Interfaces;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Infrastructure.Data.Proxies
{
    public class OperationRepositoryProxy : IOperationRepository
    {
        private readonly IOperationRepository _repository;
        private readonly Dictionary<Guid, Operation> _cache = new();

        public OperationRepositoryProxy(IOperationRepository repository)
        {
            _repository = repository;

            var allOperations = _repository.GetAllOperations();
            foreach (var operation in allOperations)
            {
                _cache[operation.Id] = operation;
            }
        }

        public Operation CreateOperation(ItemType type, Guid bankAccountId, decimal amount, DateTime date, Guid categoryId, string? description = null)
        {
            var operation = _repository.CreateOperation(type, bankAccountId, amount, date, categoryId, description);

            _cache[operation.Id] = operation;

            return operation;
        }

        public Operation? GetOperation(Guid operationId)
        {
            if (_cache.TryGetValue(operationId, out var operation))
                return operation;

            operation = _repository.GetOperation(operationId);
            if (operation != null)
            {
                _cache[operation.Id] = operation;
            }

            return operation;
        }

        public IEnumerable<Operation> GetAllOperations()
        {
            return _cache.Values;
        }

        public bool DeleteOperation(Guid operationId)
        {
            var success = _repository.DeleteOperation(operationId);

            if (success)
            {
                _cache.Remove(operationId);
            }

            return success;
        }

        public void UpdateOperation(Operation operation)
        {
            _repository.UpdateOperation(operation);
            _cache[operation.Id] = operation;
        }

        public void UploadOperation(Operation operation)
        {
            if (_cache.ContainsKey(operation.Id))
            {
                throw new InvalidOperationException($"An operation with ID {operation.Id} already exists.");
            }
        
            _repository.UploadOperation(operation);
            _cache[operation.Id] = operation;
        }

        public void Accept(IVisitor visitor)
        {
            foreach (var account in GetAllOperations())
            {
                account.Accept(visitor);
            }
        }
    }
}