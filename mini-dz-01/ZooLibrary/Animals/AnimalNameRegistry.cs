using System;
using System.Collections.Generic;
using ZooLibrary.Animals;

namespace ZooLibrary.Animals
{
    /// <summary>
    /// Реестр для хранения, получения названий животных и фабричного создания.
    /// </summary>
    public class AnimalNameRegistry
    {
        private static readonly Lazy<AnimalNameRegistry> _instance = new(() => new AnimalNameRegistry());
        
        private readonly Dictionary<Type, string> _animalNames = new();
        private readonly Dictionary<string, Type> _nameToAnimalType = new();
        
        private AnimalNameRegistry()
        {
            // Инициализация при первом обращении к реестру
            RegisterAllAnimals();
        }
        
        public static AnimalNameRegistry Instance => _instance.Value;

        public void RegisterAnimalName(Type animalType, string name)
        {
            if (!_animalNames.ContainsKey(animalType))
            {
                _animalNames[animalType] = name;
                _nameToAnimalType[name] = animalType;
            }
        }

        public string GetAnimalName<TAnimal>() where TAnimal : Animal
        {
            var animalType = typeof(TAnimal);
            return GetAnimalName(animalType);
        }

        public string GetAnimalName(Type animalType)
        {
            return _animalNames.TryGetValue(animalType, out var name) ? name : "Неизвестное животное";
        }

        public Type? GetAnimalType(string name)
        {
            return _nameToAnimalType.TryGetValue(name, out var animalType) ? animalType : null;
        }

        public IEnumerable<string> GetAnimalNames()
        {
            return _nameToAnimalType.Keys;
        }
        
        /// <summary>
        /// Регистрирует всех животных, найденных в текущей сборке.
        /// </summary>
        private void RegisterAllAnimals()
        {
            var animalTypes = AppDomain.CurrentDomain.GetAssemblies()
                .SelectMany(assembly => assembly.GetTypes())
                .Where(type => typeof(Animal).IsAssignableFrom(type) && !type.IsAbstract);

            foreach (var type in animalTypes)
            {
                var attribute = type.GetCustomAttributes(typeof(AnimalNameAttribute), false)
                    .Cast<AnimalNameAttribute>()
                    .FirstOrDefault();

                if (attribute != null)
                {
                    RegisterAnimalName(type, attribute.Name);
                }
                else
                {
                    RegisterAnimalName(type, "Неизвестное животное");
                }
            }

        }
    }
}