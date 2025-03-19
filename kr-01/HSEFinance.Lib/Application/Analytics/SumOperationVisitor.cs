using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Application.Analytics
{
    public class SumOperationVisitor : IVisitor
    {
        public decimal TotalIncome { get; private set; }
        public decimal TotalExpense { get; private set; }

        public void Visit(BankAccount bankAccount)
        {
            // Аналитика для счетов не требуется в данном посетителе
        }

        public void Visit(Category category)
        {
            // Аналитика для категорий не требуется в данном посетителе
        }

        public void Visit(Operation operation)
        {
            if (operation.Type == ItemType.Income)
            {
                TotalIncome += operation.Amount;
            }
            else if (operation.Type == ItemType.Expense)
            {
                TotalExpense += operation.Amount;
            }
        }

        public (decimal Income, decimal Expense) GetResults()
        {
            return (TotalIncome, TotalExpense);
        }
    }
}