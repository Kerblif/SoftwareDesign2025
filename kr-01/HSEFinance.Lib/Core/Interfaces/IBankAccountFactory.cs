using HSEFinance.Lib.Domain.Entities;

namespace HSEFinance.Lib.Core
{
    public interface IBankAccountFactory
    {
        BankAccount Create(string name);
    }
}