using HSEFinance.Lib.Domain.Entities;

namespace HSEFinance.Lib.Domain.Repositories
{
    public interface IAccountRepository
    {
        IEnumerable<BankAccount> GetAllBankAccounts();
        BankAccount CreateBankAccount(string? name);
        BankAccount? GetBankAccount(Guid accountId);
        bool DeleteBankAccount(Guid accountId);
        void UpdateBankAccount(BankAccount account);
        void UploadBankAccount(BankAccount account);
        void RecalculateAccountBalance(Guid accountId);
    }
}