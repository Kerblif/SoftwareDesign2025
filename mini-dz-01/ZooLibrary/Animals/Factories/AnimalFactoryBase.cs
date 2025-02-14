using System;

namespace ZooLibrary.Animals.Factories
{
    public abstract class AnimalFactoryBase<TAnimal> : IAutoRegisteredAnimalFactory where TAnimal : Animal
    { 
        public Type GetAnimalType() => typeof(TAnimal);

        public abstract Animal CreateAnimal(int food, int number, bool isHealthy);
    }
}