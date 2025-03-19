// Файл: HSEFinance.Lib/Infrastructure/Data/OperationRepository.cs

using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using System;
using System.Collections.Generic;
using System.Linq;

namespace HSEFinance.Lib.Infrastructure.Data
{
    public class OperationRepository : IOperationRepository
    {
        private readonly HSEFinanceDbContext _dbContext;
        
        public OperationRepository(HSEFinanceDbContext dbContext)
        {
            _dbContext = dbContext;
        }

        public IEnumerable<Operation> GetAllOperations()
        {
            return _dbContext.Operations;
        }

        public Operation CreateOperation(ItemType type,
                                           Guid bankAccountId,
                                           decimal amount,
                                           DateTime date,
                                           Guid categoryId,
                                           string? description = null)
        {
            var operation = new Operation(type, bankAccountId, amount, date, categoryId, description);
            _dbContext.Operations.Add(operation);
            _dbContext.SaveChanges();
            return operation;
        }

        public Operation? GetOperation(Guid operationId)
        {
            return _dbContext.Operations.Find(operationId);
        }

        public bool DeleteOperation(Guid operationId)
        {
            var operation = GetOperation(operationId);
            if (operation == null)
                return false;

            _dbContext.Operations.Remove(operation);
            _dbContext.SaveChanges();
            return true;
        }
        
        public void UpdateOperation(Operation operation)
        {
            _dbContext.Operations.Update(operation);
            _dbContext.SaveChanges();
        }
    }
}