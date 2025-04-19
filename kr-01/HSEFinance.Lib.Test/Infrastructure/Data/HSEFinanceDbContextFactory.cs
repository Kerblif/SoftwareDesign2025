// Тестовая база данных: HSEFinanceDbContextFactory
using HSEFinance.Lib.Infrastructure.Data;
using Microsoft.EntityFrameworkCore;

namespace HSEFinance.Lib.Test.Infrastructure
{
    public static class HSEFinanceDbContextFactory
    {
        private static void Clear(HSEFinanceDbContext context)
        {
            foreach (var entity in context.BankAccounts)
            {
                context.BankAccounts.Remove(entity);
            }

            foreach (var entity in context.Categories)
            {
                context.Categories.Remove(entity);
            }

            foreach (var entity in context.Operations)
            {
                context.Operations.Remove(entity);
            }
        }
        
        public static HSEFinanceDbContext Create()
        {
            var optionsBuilder = new DbContextOptionsBuilder<HSEFinanceDbContext>();
        
            optionsBuilder.UseSqlite("Data Source=/Users/kerblif/Программирование/HSE/KPO/kr-01/app.db");
        
            var dbContext = new HSEFinanceDbContext(optionsBuilder.Options);
            Clear(dbContext);
            
            return dbContext;
        }
    }
}