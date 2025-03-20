// Файл: HSEFinance.Lib/Infrastructure/Data/OperationRepository.cs

using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Domain.Repositories;
using System;
using System.Collections.Generic;
using System.Linq;

namespace HSEFinance.Lib.Infrastructure.Data
{
    public class OperationRepository : IOperationRepository
    {
        private readonly HSEFinanceDbContext _dbContext;
        private readonly IAccountRepository _accountRepository;
        private readonly ICategoryRepository _categoryRepository;
        
        public OperationRepository(HSEFinanceDbContext dbContext, IAccountRepository accountRepository, ICategoryRepository categoryRepository)
        {
            _dbContext = dbContext;
            _accountRepository = accountRepository;
            _categoryRepository = categoryRepository;
        }

        public IEnumerable<Operation> GetAllOperations()
        {
            return _dbContext.Operations;
        }

        public Operation CreateOperation(ItemType type,
                                           Guid bankAccountId,
                                           decimal amount,
                                           DateTime date,
                                           Guid categoryId,
                                           string? description = null)
        {
            using var transaction = _dbContext.Database.BeginTransaction();
            
            try
            {
                BankAccount? bankAccount = _accountRepository.GetBankAccount(bankAccountId);
                Category? category = _categoryRepository.GetCategory(categoryId);

                if (bankAccount == null)
                {
                    throw new Exception("Bank account not found");
                }

                if (category == null)
                {
                    throw new Exception("Category not found");
                }

                var operation = new Operation(type, bankAccountId, amount, date, categoryId, description);
                _dbContext.Operations.Add(operation);
                _dbContext.SaveChanges();
                
                bankAccount.UpdateBalance(amount, type);
                _accountRepository.UpdateBankAccount(bankAccount);

                transaction.Commit();

                return operation;
            }
            catch (Exception)
            {
                transaction.Rollback();
                throw;
            }
        }

        public Operation? GetOperation(Guid operationId)
        {
            return _dbContext.Operations.Find(operationId);
        }

        public bool DeleteOperation(Guid operationId)
        {
            Operation? operation = GetOperation(operationId);

            if (operation == null)
            {
                throw new Exception("Operation not found");
            }
            
            using var transaction = _dbContext.Database.BeginTransaction();
            
            try
            {
                BankAccount? bankAccount = _accountRepository.GetBankAccount(operation.BankAccountId);

                if (bankAccount == null)
                {
                    throw new Exception("Bank account not found");
                }

                _dbContext.Operations.Remove(operation);
                _dbContext.SaveChanges();
                
                bankAccount.UpdateBalance(operation.Amount, operation.Type == ItemType.Income ? ItemType.Expense : ItemType.Income);
                _accountRepository.UpdateBankAccount(bankAccount);

                transaction.Commit();

                return true;
            }
            catch (Exception)
            {
                transaction.Rollback();
                throw;
            }
        }
        
        public void UpdateOperation(Operation operation)
        {
            var existingOperation = _dbContext.Operations.Find(operation.Id);
            if (existingOperation == null)
            {
                throw new Exception("Operation not found.");
            }
        
            if (existingOperation.Type != operation.Type ||
                existingOperation.BankAccountId != operation.BankAccountId ||
                existingOperation.Amount != operation.Amount ||
                existingOperation.Date != operation.Date)
            {
                throw new InvalidOperationException("Only the comment can be updated.");
            }

            if (_categoryRepository.GetCategory(operation.CategoryId) == null)
            {
                throw new Exception("Category not found.");
            }
        
            existingOperation.CategoryId = operation.CategoryId;
            existingOperation.Description = operation.Description;
            _dbContext.Operations.Update(existingOperation);
            _dbContext.SaveChanges();
        }
    }
}