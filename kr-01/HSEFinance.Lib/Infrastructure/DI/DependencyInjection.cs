using HSEFinance.Lib.Application.Facades;
using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Factories;
using HSEFinance.Lib.Domain.Repositories;
using HSEFinance.Lib.Infrastructure.Data;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;

namespace HSEFinance.Lib.Infrastructure.DI
{
    public static class DependencyInjection
    {
        public static IServiceCollection AddHSEFinanceServices(this IServiceCollection services, string connectionString)
        {
            // Настройка DbContext для SQLite
            services.AddDbContext<HSEFinanceDbContext>(options =>
                options.UseSqlite(connectionString));

            // Регистрация фабрик
            services.AddScoped<IBankAccountFactory, BankAccountFactory>();
            services.AddScoped<ICategoryFactory, CategoryFactory>();
            services.AddScoped<IOperationFactory, OperationFactory>();

            // Регистрация репозиториев
            services.AddScoped<IAccountRepository, AccountRepository>();
            services.AddScoped<ICategoryRepository, CategoryRepository>();
            services.AddScoped<IOperationRepository, OperationRepository>();

            // Регистрация фасадов
            services.AddScoped<BankAccountFacade>();
            services.AddScoped<CategoryFacade>();
            services.AddScoped<OperationFacade>();

            return services;
        }
    }
}