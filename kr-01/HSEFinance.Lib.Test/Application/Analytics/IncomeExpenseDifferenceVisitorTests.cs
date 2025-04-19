using HSEFinance.Lib.Application.Analytics;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Analytics
{
    public class IncomeExpenseDifferenceVisitorTests
    {
        [Fact]
        public void Visit_WithIncomeAndExpenseInRange_CalculatesDifference()
        {
            // Arrange
            var startDate = new DateTime(2023, 1, 1);
            var endDate = new DateTime(2023, 12, 31);
            var visitor = new IncomeExpenseDifferenceVisitor(startDate, endDate);

            var incomeOperation = new Operation(ItemType.Income, Guid.NewGuid(),  100m, new DateTime(2023, 1, 1), Guid.NewGuid());
            var expenseOperation = new Operation(ItemType.Expense, Guid.NewGuid(),  50m, new DateTime(2023, 1, 1), Guid.NewGuid());

            // Act
            visitor.Visit(incomeOperation);
            visitor.Visit(expenseOperation);

            // Assert
            Assert.Equal(100m, visitor.TotalIncome);
            Assert.Equal(50m, visitor.TotalExpense);
            Assert.Equal(50m, visitor.CalculateDifference());
        }

        [Fact]
        public void Visit_WithOperationsOutsideDateRange_IgnoresThem()
        {
            // Arrange
            var startDate = new DateTime(2023, 1, 1);
            var endDate = new DateTime(2023, 12, 31);
            var visitor = new IncomeExpenseDifferenceVisitor(startDate, endDate);

            var incomeOperation = new Operation(ItemType.Income, Guid.NewGuid(),  100m, new DateTime(2022, 12, 31), Guid.NewGuid());
            var expenseOperation = new Operation(ItemType.Expense, Guid.NewGuid(),  50m, new DateTime(2024, 2, 1), Guid.NewGuid());

            // Act
            visitor.Visit(incomeOperation);
            visitor.Visit(expenseOperation);

            // Assert
            Assert.Equal(0m, visitor.TotalIncome);
            Assert.Equal(0m, visitor.TotalExpense);
            Assert.Equal(0m, visitor.CalculateDifference());
        }
    }
}