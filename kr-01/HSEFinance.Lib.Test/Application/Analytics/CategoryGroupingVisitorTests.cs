using HSEFinance.Lib.Application.Analytics;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using Xunit;

namespace HSEFinance.Lib.Test.Application.Analytics
{
    public class CategoryGroupingVisitorTests
    {
        [Fact]
        public void Visit_WithIncomeAndExpenseOperations_GroupsAmountsByCategory()
        {
            var categoryId1 = Guid.NewGuid();
            var categoryId2 = Guid.NewGuid();
            var visitor = new CategoryGroupingVisitor();

            var incomeOperation1 = new Operation(ItemType.Income, Guid.NewGuid(),  100m, DateTime.Now, categoryId1);
            var incomeOperation2 = new Operation(ItemType.Income, Guid.NewGuid(),  200m, DateTime.Now, categoryId1);
            var expenseOperation = new Operation(ItemType.Expense, Guid.NewGuid(),  50m, DateTime.Now, categoryId2);

            visitor.Visit(incomeOperation1);
            visitor.Visit(incomeOperation2);
            visitor.Visit(expenseOperation);

            Assert.Equal(300m, visitor.IncomeByCategory[categoryId1]);
            Assert.Equal(50m, visitor.ExpenseByCategory[categoryId2]);
        }

        [Fact]
        public void Visit_WithNoOperations_ProduceEmptyDictionaries()
        {
            var visitor = new CategoryGroupingVisitor();

            Assert.Empty(visitor.IncomeByCategory);
            Assert.Empty(visitor.ExpenseByCategory);
        }
    }
}