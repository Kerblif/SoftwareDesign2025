using System.Linq;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Domain.Factories;
using HSEFinance.Lib.Infrastructure.Data;
using Xunit;

namespace HSEFinance.Lib.Test.Infrastructure.Data
{
    public class CategoryRepositoryTests
    {
        [Fact]
        public void CreateCategory_ShouldAddCategoryToDatabase()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var categoryFactory = new CategoryFactory();
            var repository = new CategoryRepository(dbContext, categoryFactory);

            // Act
            var category = repository.CreateCategory(ItemType.Expense, "Groceries");

            // Assert
            Assert.NotNull(category);
            Assert.Equal("Groceries", category.Name);
            Assert.Single(dbContext.Categories);
        }

        [Fact]
        public void GetCategory_WithValidId_ReturnsCategory()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var categoryFactory = new CategoryFactory();
            var repository = new CategoryRepository(dbContext, categoryFactory);

            var createdCategory = repository.CreateCategory(ItemType.Income, "Salary");

            // Act
            var retrievedCategory = repository.GetCategory(createdCategory.Id);

            // Assert
            Assert.NotNull(retrievedCategory);
            Assert.Equal("Salary", retrievedCategory?.Name);
        }

        [Fact]
        public void DeleteCategory_RemovesCategoryFromDatabase()
        {
            // Arrange
            var dbContext = HSEFinanceDbContextFactory.Create();
            var categoryFactory = new CategoryFactory();
            var repository = new CategoryRepository(dbContext, categoryFactory);

            var category = repository.CreateCategory(ItemType.Expense, "Travel");

            // Act
            var result = repository.DeleteCategory(category.Id);

            // Assert
            Assert.True(result);
            Assert.Empty(dbContext.Categories);
        }
    }
}