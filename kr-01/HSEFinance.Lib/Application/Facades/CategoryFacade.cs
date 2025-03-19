using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;

namespace HSEFinance.Lib.Application.Facades
{
    public class CategoryFacade
    {
        private readonly ICategoryRepository _categoryRepository;
        private readonly ICategoryFactory _categoryFactory;

        public CategoryFacade(ICategoryRepository categoryRepository, ICategoryFactory categoryFactory)
        {
            _categoryRepository = categoryRepository ?? throw new ArgumentNullException(nameof(categoryRepository));
            _categoryFactory = categoryFactory ?? throw new ArgumentNullException(nameof(categoryFactory));
        }

        public Category CreateCategory(ItemType type, string name)
        {
            var category = _categoryFactory.Create(type, name);

            _categoryRepository.CreateCategory(type, name);

            return category;
        }

        public Category? GetCategory(Guid categoryId)
        {
            return _categoryRepository.GetCategory(categoryId);
        }

        public bool DeleteCategory(Guid categoryId)
        {
            return _categoryRepository.DeleteCategory(categoryId);
        }
    }
}