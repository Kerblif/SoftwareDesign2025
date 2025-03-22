using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data.Proxies;
using Moq;
using Xunit;

namespace HSEFinance.Lib.Test.Infrastructure.Data.Proxies
{
    public class CategoryRepositoryProxyTests
    {
        [Fact]
        public void GetAllCategories_InitiallyLoadsAllCategoriesFromBaseRepository()
        {
            // Arrange
            var mockRepository = new Mock<ICategoryRepository>();
            var categories = new List<Category>
            {
                new Category(ItemType.Expense, "Food"),
                new Category(ItemType.Income, "Salary")
            };
            mockRepository.Setup(r => r.GetAllCategories()).Returns(categories);

            var proxy = new CategoryRepositoryProxy(mockRepository.Object);

            // Act
            var allCategories = proxy.GetAllCategories();

            // Assert
            Assert.Equal(2, allCategories.Count());
            mockRepository.Verify(r => r.GetAllCategories(), Times.Once);
        }

        [Fact]
        public void GetCategory_UsesCacheForPreviouslyLoadedCategory()
        {
            // Arrange
            var mockRepository = new Mock<ICategoryRepository>();
            var category = new Category(ItemType.Income, "Cached Category");
            mockRepository.Setup(r => r.GetAllCategories()).Returns(new List<Category> { category });

            var proxy = new CategoryRepositoryProxy(mockRepository.Object);

            // Act
            var retrievedCategory = proxy.GetCategory(category.Id);

            // Assert
            Assert.NotNull(retrievedCategory);
            Assert.Equal(category.Name, retrievedCategory!.Name);
            mockRepository.Verify(r => r.GetCategory(It.IsAny<Guid>()), Times.Never);
        }

        [Fact]
        public void CreateCategory_AddsToRepositoryAndCache()
        {
            // Arrange
            var mockRepository = new Mock<ICategoryRepository>();
            var category = new Category(ItemType.Expense, "New Category");
            mockRepository.Setup(r => r.CreateCategory(ItemType.Expense, "New Category")).Returns(category);

            var proxy = new CategoryRepositoryProxy(mockRepository.Object);

            // Act
            var createdCategory = proxy.CreateCategory(ItemType.Expense, "New Category");

            // Assert
            Assert.NotNull(createdCategory);
            Assert.Equal("New Category", createdCategory.Name);
            mockRepository.Verify(r => r.CreateCategory(ItemType.Expense, "New Category"), Times.Once);

            var cachedCategory = proxy.GetCategory(category.Id);
            Assert.NotNull(cachedCategory);
            Assert.Equal("New Category", cachedCategory!.Name);
        }

        [Fact]
        public void DeleteCategory_RemovesFromCacheAndRepository()
        {
            // Arrange
            var mockRepository = new Mock<ICategoryRepository>();
            var categories = new List<Category>
            {
                new Category(ItemType.Income, "To Be Deleted")
            };
            mockRepository.Setup(r => r.GetAllCategories()).Returns(categories);
            mockRepository.Setup(r => r.DeleteCategory(categories[0].Id)).Returns(true);

            var proxy = new CategoryRepositoryProxy(mockRepository.Object);

            // Act
            var deleteResult = proxy.DeleteCategory(categories[0].Id);

            // Assert
            Assert.True(deleteResult);
            Assert.Null(proxy.GetCategory(categories[0].Id));
            mockRepository.Verify(r => r.DeleteCategory(categories[0].Id), Times.Once);
        }
    }
}