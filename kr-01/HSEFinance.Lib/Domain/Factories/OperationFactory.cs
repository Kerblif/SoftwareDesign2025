using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using System;

namespace HSEFinance.Lib.Domain.Factories
{
    public class OperationFactory : IOperationFactory
    {
        public Operation Create(
            ItemType type, 
            Guid bankAccountId, 
            decimal amount, 
            DateTime date, 
            Guid categoryId, 
            string? description = null)
        {
            return new Operation(type, bankAccountId, amount, date, categoryId, description);
        }
    }
}