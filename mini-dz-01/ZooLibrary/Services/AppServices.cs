using Microsoft.Extensions.DependencyInjection;
using ZooLibrary.Animals;
using ZooLibrary.Animals.Factories;
using ZooLibrary.ZooEntities;

namespace ZooLibrary.Services
{
    public static class AppServices
    {
        private static IServiceProvider? _services;
        private static readonly object _lockObj = new();

        /// <summary>
        /// Провайдер сервисов приложения.
        /// </summary>
        public static IServiceProvider Services
        {
            get
            {
                lock (_lockObj)
                {
                    _services ??= ConfigureServices();
                    return _services;
                }
            }
        }

        /// <summary>
        /// Инициализация и настройка DI-контейнера.
        /// </summary>
        private static IServiceProvider ConfigureServices()
        {
            Console.WriteLine("Настройка сервисов...");

            var services = new ServiceCollection();

            services.AddSingleton(AnimalFactoryRegistry.Instance);
            services.AddSingleton(AnimalNameRegistry.Instance);
            services.AddSingleton(Zoo.Instance);

            return services.BuildServiceProvider();
        }
    }
}