using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;

namespace HSEFinance.Lib.Domain.Factories
{
    public class BankAccountFactory : IBankAccountFactory
    {
        public BankAccount Create(string? name)
        {
            if (string.IsNullOrWhiteSpace(name))
            {
                throw new ArgumentException("The account name cannot be empty.");
            }
            
            return new BankAccount(name);
        }
    }
}