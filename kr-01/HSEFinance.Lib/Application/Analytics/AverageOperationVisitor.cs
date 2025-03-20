using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Application.Analytics
{
    public class AverageOperationVisitor : IVisitor
    {
        private decimal _totalIncome;
        private decimal _totalExpense;
        private int _incomeCount;
        private int _expenseCount;

        public void Visit(BankAccount bankAccount)
        {
            // Анализировать счета не требуется
        }

        public void Visit(Category category)
        {
            // Анализировать категории не требуется
        }

        public void Visit(Operation operation)
        {
            if (operation.Type == ItemType.Income)
            {
                _totalIncome += operation.Amount;
                _incomeCount++;
            }
            else if (operation.Type == ItemType.Expense)
            {
                _totalExpense += operation.Amount;
                _expenseCount++;
            }
        }

        public decimal GetAverageIncome()
        {
            return _incomeCount > 0 ? _totalIncome / _incomeCount : 0;
        }

        public decimal GetAverageExpense()
        {
            return _expenseCount > 0 ? _totalExpense / _expenseCount : 0;
        }
    }
}