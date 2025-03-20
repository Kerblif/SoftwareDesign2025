using System.Collections.Generic;
using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Application.Analytics
{
    public class CategoryGroupingVisitor : IVisitor
    {
        public Dictionary<Guid, decimal> IncomeByCategory { get; } = new();
        public Dictionary<Guid, decimal> ExpenseByCategory { get; } = new();

        public void Visit(BankAccount bankAccount)
        {
            // Не требуется анализировать счета
        }

        public void Visit(Category category)
        {
            // Не требуется анализировать категории напрямую
        }

        public void Visit(Operation operation)
        {
            if (operation.Type == ItemType.Income)
            {
                if (!IncomeByCategory.ContainsKey(operation.CategoryId))
                {
                    IncomeByCategory[operation.CategoryId] = 0;
                }
                IncomeByCategory[operation.CategoryId] += operation.Amount;
            }
            else if (operation.Type == ItemType.Expense)
            {
                if (!ExpenseByCategory.ContainsKey(operation.CategoryId))
                {
                    ExpenseByCategory[operation.CategoryId] = 0;
                }
                ExpenseByCategory[operation.CategoryId] += operation.Amount;
            }
        }
    }
}