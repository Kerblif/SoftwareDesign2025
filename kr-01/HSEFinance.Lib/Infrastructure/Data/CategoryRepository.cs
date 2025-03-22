using HSEFinance.Lib.Core;
using HSEFinance.Lib.Core.Interfaces;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using Microsoft.EntityFrameworkCore;

namespace HSEFinance.Lib.Infrastructure.Data
{
    public class CategoryRepository : ICategoryRepository
    {
        private readonly HSEFinanceDbContext _dbContext;

        public CategoryRepository(HSEFinanceDbContext dbContext)
        {
            _dbContext = dbContext;
        }

        public IEnumerable<Category> GetAllCategories()
        {
            return _dbContext.Categories;
        }

        public Category CreateCategory(ItemType type, string name)
        {
            var category = new Category(type, name);
            _dbContext.Categories.Add(category);
            _dbContext.SaveChanges();
            return category;
        }

        public Category? GetCategory(Guid categoryId)
        {
            return _dbContext.Categories.Find(categoryId);
        }

        public bool DeleteCategory(Guid categoryId)
        {
            var category = GetCategory(categoryId);
            if (category == null)
                return false;

            _dbContext.Categories.Remove(category);
            _dbContext.SaveChanges();
            return true;
        }

        public void UpdateCategory(Category category)
        {
            _dbContext.Categories.Update(category);
            _dbContext.SaveChanges();
        }

        public void UploadCategory(Category account)
        {
            if (_dbContext.Categories.Find(account.Id) != null)
            {
                throw new InvalidOperationException($"A category with ID {account.Id} already exists.");
            }
        
            _dbContext.Categories.Add(account);
            _dbContext.SaveChanges();
        }

        public void Accept(IVisitor visitor)
        {
            foreach (var account in GetAllCategories())
            {
                account.Accept(visitor);
            }
        }
    }
}