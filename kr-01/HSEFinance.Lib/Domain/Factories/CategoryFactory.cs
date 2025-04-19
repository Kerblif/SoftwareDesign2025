using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Domain.Factories
{
    public class CategoryFactory : ICategoryFactory
    {
        public Category Create(ItemType type, string name)
        {
            if (string.IsNullOrWhiteSpace(name))
            {
                throw new ArgumentException("Category name cannot be empty.");
            }
            
            return new Category(type, name);
        }
    }
}