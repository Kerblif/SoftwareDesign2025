using HSEFinance.Lib.Domain.Entities;
using Microsoft.EntityFrameworkCore;

namespace HSEFinance.Lib.Infrastructure.Data
{
    public class HSEFinanceDbContext : DbContext
    {
        public DbSet<BankAccount> BankAccounts { get; set; } = null!;
        public DbSet<Category> Categories { get; set; } = null!;
        public DbSet<Operation> Operations { get; set; } = null!;

        public HSEFinanceDbContext(DbContextOptions<HSEFinanceDbContext> options)
            : base(options)
        {
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<BankAccount>(entity =>
            {
                entity.HasKey(b => b.Id);
                entity.Property(b => b.Name).IsRequired();
                entity.Property(b => b.Balance).HasDefaultValue(0);
            });

            modelBuilder.Entity<Category>(entity =>
            {
                entity.HasKey(c => c.Id);
                entity.Property(c => c.Type).IsRequired();
                entity.Property(c => c.Name).IsRequired();
            });

            modelBuilder.Entity<Operation>(entity =>
            {
                entity.HasKey(o => o.Id);
                entity.Property(o => o.Type).IsRequired();
                entity.Property(o => o.BankAccountId).IsRequired();
                entity.Property(o => o.Amount).IsRequired();
                entity.Property(o => o.Date).IsRequired();
                entity.Property(o => o.Description).HasDefaultValue("");
                entity.Property(o => o.CategoryId).IsRequired();
            });

            base.OnModelCreating(modelBuilder);
        }
    }
}