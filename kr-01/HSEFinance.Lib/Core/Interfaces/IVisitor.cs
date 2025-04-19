using HSEFinance.Lib.Domain.Entities;

namespace HSEFinance.Lib.Core
{
    public interface IVisitor
    {
        void Visit(BankAccount bankAccount);
        void Visit(Category category);
        void Visit(Operation operation);
    }
}