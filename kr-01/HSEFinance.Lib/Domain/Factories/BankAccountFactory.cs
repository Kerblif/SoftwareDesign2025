using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;

namespace HSEFinance.Lib.Domain.Factories
{
    public class BankAccountFactory : IBankAccountFactory
    {
        public BankAccount Create(string name)
        {
            return new BankAccount(name);
        }
    }
}