using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Threading.Tasks;
using HSEFinance.Lib.Application.Facades;
using Microsoft.Extensions.DependencyInjection;
using Spectre.Console;
using HSEFinance.Lib.Core;
using HSEFinance.Lib.Domain.Entities;
using HSEFinance.Lib.Domain.Enums;
using HSEFinance.Lib.Infrastructure.DI;

namespace HSEFinance.ConsoleApp
{
    class Program
    {
        static void Main(string[] args)
        {
            // Настройка DI контейнера
            var serviceProvider = ConfigureServices();

            using var scope = serviceProvider.CreateScope();
            var app = scope.ServiceProvider.GetRequiredService<FinanceApp>();
            
            app.Run();
        }

        private static IServiceProvider ConfigureServices()
        {
            var services = new ServiceCollection();

            DependencyInjection.AddHSEFinanceServices(services, "Data Source=" + Environment.GetEnvironmentVariable("DATABASE_PATH"));
                
            // Регистрация фасадов
            services.AddSingleton<AccountManagerFacade>();
            services.AddSingleton<CategoryManagerFacade>();
            services.AddSingleton<OperationManagerFacade>();
            services.AddSingleton<AnalyticsFacade>();
            
            // Регистрация приложения
            services.AddSingleton<FinanceApp>();
            
            return services.BuildServiceProvider();
        }
    }
}