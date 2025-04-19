using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Design;

namespace HSEFinance.Lib.Infrastructure.Data
{
    public class HSEFinanceDbContextFactory : IDesignTimeDbContextFactory<HSEFinanceDbContext>
    {
        public HSEFinanceDbContext CreateDbContext(string[] args)
        {
            var optionsBuilder = new DbContextOptionsBuilder<HSEFinanceDbContext>();

            optionsBuilder.UseSqlite("Data Source=" + Environment.GetEnvironmentVariable("DATABASE_PATH"));

            return new HSEFinanceDbContext(optionsBuilder.Options);
        }
    }
}