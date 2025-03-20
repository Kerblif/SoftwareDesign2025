using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Infrastructure.Data.Proxies
{
    public class CategoryRepositoryProxy : ICategoryRepository
    {
        private readonly ICategoryRepository _repository;
        private readonly Dictionary<Guid, Category> _cache = new();

        public CategoryRepositoryProxy(ICategoryRepository repository)
        {
            _repository = repository;

            var allCategories = _repository.GetAllCategories();
            foreach (var category in allCategories)
            {
                _cache[category.Id] = category;
            }
        }

        public Category CreateCategory(ItemType type, string name)
        {
            var category = _repository.CreateCategory(type, name);

            _cache[category.Id] = category;

            return category;
        }

        public Category? GetCategory(Guid categoryId)
        {
            if (_cache.TryGetValue(categoryId, out var category))
                return category;

            category = _repository.GetCategory(categoryId);
            if (category != null)
            {
                _cache[category.Id] = category;
            }

            return category;
        }

        public IEnumerable<Category> GetAllCategories()
        {
            return _cache.Values;
        }

        public bool DeleteCategory(Guid categoryId)
        {
            var success = _repository.DeleteCategory(categoryId);

            if (success)
            {
                _cache.Remove(categoryId);
            }

            return success;
        }

        public void UpdateCategory(Category category)
        {
            _repository.UpdateCategory(category);
            _cache[category.Id] = category;
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