using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using System;

namespace HSEFinance.Lib.Domain.Repositories
{
    public interface ICategoryRepository
    {
        IEnumerable<Category> GetAllCategories();
        Category CreateCategory(ItemType type, string name);
        Category? GetCategory(Guid categoryId);
        bool DeleteCategory(Guid categoryId);
        void UpdateCategory(Category account);
    }
}
