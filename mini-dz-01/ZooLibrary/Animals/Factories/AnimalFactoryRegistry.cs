using System;
using System.Collections.Generic;
using System.Linq;
using Microsoft.Extensions.DependencyInjection;
using ZooLibrary.Services;

namespace ZooLibrary.Animals.Factories
{
    public class AnimalFactoryRegistry
    {
        private static readonly Lazy<AnimalFactoryRegistry> _instance = new(() => new AnimalFactoryRegistry());

        private readonly Dictionary<Type, IAutoRegisteredAnimalFactory> _factories = new();

        private AnimalFactoryRegistry()
        {
            // Инициализация при первом обращении к реестру
            RegisterAllFactories();
        }

        public static AnimalFactoryRegistry Instance => _instance.Value;

        /// <summary>
        /// Регистрирует конкретную фабрику для указанного животного.
        /// </summary>
        public void RegisterFactory<TAnimal>(IAutoRegisteredAnimalFactory factory) where TAnimal : Animal
        {
            _factories[typeof(TAnimal)] = factory;
        }

        /// <summary>
        /// Получить фабрику для указанного животного.
        /// </summary>
        public IAutoRegisteredAnimalFactory? GetFactory<TAnimal>() where TAnimal : Animal
        {
            var animalType = typeof(TAnimal);
            return _factories.TryGetValue(animalType, out var factory) ? factory : null;
        }
        
        /// <summary>
        /// Получить фабрику для указанного животного.
        /// </summary>
        public IAutoRegisteredAnimalFactory? GetFactory(string animalName)
        {
            var animalType = AppServices.Services.GetRequiredService<AnimalNameRegistry>().GetAnimalType(animalName);

            if (animalType == null || !typeof(Animal).IsAssignableFrom(animalType))
            {
                return null;
            }
            
            return _factories.GetValueOrDefault(animalType);
        }

        /// <summary>
        /// Получить список всех зарегистрированных видов.
        /// </summary>
        public IEnumerable<Type> GetRegisteredSpecies()
        {
            return _factories.Keys;
        }

        /// <summary>
        /// Регистрирует все фабрики, найденные в текущей сборке.
        /// </summary>
        private void RegisterAllFactories()
        {
            var factoryTypes = AppDomain.CurrentDomain.GetAssemblies()
                .SelectMany(assembly => assembly.GetTypes())
                .Where(t => typeof(IAutoRegisteredAnimalFactory).IsAssignableFrom(t) && !t.IsInterface && !t.IsAbstract);

            foreach (var factoryType in factoryTypes)
            {
                var factoryInstance = (IAutoRegisteredAnimalFactory)Activator.CreateInstance(factoryType)!;
                var animalType = factoryInstance.GetAnimalType();
                _factories[animalType] = factoryInstance;
            }
        }
    }
}