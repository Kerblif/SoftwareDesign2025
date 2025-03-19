using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;

namespace HSEFinance.Lib.Core
{
    public interface ICategoryFactory
    {
        Category Create(ItemType type, string name);
    }
}