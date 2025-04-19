using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using System;

namespace HSEFinance.Lib.Core
{
    public interface IOperationFactory
    {
        Operation Create(
            ItemType type, 
            Guid bankAccountId, 
            decimal amount, 
            DateTime date, 
            Guid categoryId, 
            string? description = null);
    }
}