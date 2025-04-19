using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Application.Analytics
{
    public class IncomeExpenseDifferenceVisitor : IVisitor
    {
        public decimal TotalIncome { get; private set; }
        public decimal TotalExpense { get; private set; }
        public DateTime StartDate { get; }
        public DateTime EndDate { get; }
        
        public IncomeExpenseDifferenceVisitor(DateTime startDate, DateTime endDate)
        {
            StartDate = startDate;
            EndDate = endDate;
        }
        
        public void Visit(BankAccount bankAccount)
        {
            // Не требуется анализировать счета
        }

        public void Visit(Category category)
        {
            // Не требуется анализировать категории
        }

        public void Visit(Operation operation)
        {
            if (operation.Date >= StartDate && operation.Date <= EndDate)
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
        }

        public decimal CalculateDifference()
        {
            return TotalIncome - TotalExpense;
        }
    }
}