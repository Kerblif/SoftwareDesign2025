using HSEFinance.Lib.Core;
using HSEFinance.Lib.Core.Interfaces;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Interfaces;

namespace HSEFinance.Lib.Domain.Entities
{
    public class Operation : IIdentifiable, IVisitable
    {
        public Guid Id { get; }
        public ItemType Type { get; }
        public Guid BankAccountId { get; }
        public decimal Amount { get; }
        public DateTime Date { get; }
        public string? Description { get; set; }
        public Guid CategoryId { get; set; }

        public Operation(ItemType type, Guid bankAccountId, decimal amount, DateTime date, Guid categoryId, string? description = null)
        {
            Id = Guid.NewGuid();
            Type = type;
            BankAccountId = bankAccountId;
            Amount = amount;
            Date = date;
            Description = description;
            CategoryId = categoryId;

            if (amount < 0)
            {
                throw new ArgumentException("Operation amount cannot be negative.");
            }
        }
        
        public void Accept(IVisitor visitor)
        {
            visitor.Visit(this);
        }

        public override string ToString()
        {
            return $"Операция: {{ Тип: {Type}, Сумма: {Amount}, Дата: {Date}, Описание: {Description ?? "N/A"}, ID категории: {CategoryId}, ID счета: {BankAccountId} }}";
        }
    }
}