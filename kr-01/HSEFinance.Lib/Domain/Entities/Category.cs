using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Interfaces;
using HSEFinance.Lib.Core;

namespace HSEFinance.Lib.Domain.Entities
{
    public class Category : IIdentifiable, IVisitable
    {
        public Guid Id { get; }
        public ItemType Type { get; set; }
        public string Name { get; set; }

        public Category(ItemType type, string name)
        {
            Id = Guid.NewGuid();
            Type = type;
            Name = name;
        }
        
        public void Accept(IVisitor visitor)
        {
            visitor.Visit(this);
        }
    }
}