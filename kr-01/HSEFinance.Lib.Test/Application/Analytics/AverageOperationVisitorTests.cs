using HSEFinance.Lib.Application.Analytics;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Analytics
{
    public class AverageOperationVisitorTests
    {
        [Fact]
        public void Visit_WithIncomeAndExpenseOperations_CalculatesCorrectAverages()
        {
            // Arrange
            var visitor = new AverageOperationVisitor();
            var incomeOperation1 = new Operation(ItemType.Income, Guid.NewGuid(),  100m, DateTime.Now, Guid.NewGuid());
            var incomeOperation2 = new Operation(ItemType.Income, Guid.NewGuid(),  200m, DateTime.Now, Guid.NewGuid());
            var expenseOperation = new Operation(ItemType.Expense, Guid.NewGuid(),  50m, DateTime.Now, Guid.NewGuid());
            
            // Act
            visitor.Visit(incomeOperation1);
            visitor.Visit(incomeOperation2);
            visitor.Visit(expenseOperation);

            // Assert
            Assert.Equal(150m, visitor.GetAverageIncome());
            Assert.Equal(50m, visitor.GetAverageExpense());
        }

        [Fact]
        public void Visit_WithNoOperations_ReturnsZeroAverages()
        {
            // Arrange
            var visitor = new AverageOperationVisitor();

            // Act & Assert
            Assert.Equal(0m, visitor.GetAverageIncome());
            Assert.Equal(0m, visitor.GetAverageExpense());
        }
    }
}