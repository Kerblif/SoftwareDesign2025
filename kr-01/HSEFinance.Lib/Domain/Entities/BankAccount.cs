using System.Runtime.CompilerServices;
using CsvHelper.Configuration.Attributes;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Interfaces;
using HSEFinance.Lib.Core;
using HSEFinance.Lib.Core.Interfaces;

namespace HSEFinance.Lib.Domain.Entities
{
    public class BankAccount : IIdentifiable, IVisitable
    {
        public Guid Id { get; }
        
        public string? Name { get; set; }
        
        public decimal Balance { get; set; }
        
        public BankAccount(string? name)
        {
            Id = Guid.NewGuid();
            Name = name;
            Balance = 0;
        }

        public void UpdateBalance(decimal amount, ItemType type)
        {
            if (type == ItemType.Income)
            {
                Balance += amount;
            }
            else if (type == ItemType.Expense)
            {
                Balance -= amount;
            }
        }
        
        public void Accept(IVisitor visitor)
        {
            visitor.Visit(this);
        }
        
        public override string? ToString() => Name;
    }
}