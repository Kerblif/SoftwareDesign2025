using System;

namespace ZooLibrary.Animals.Factories
{
    public class RabbitFactory : AnimalFactoryBase<Rabbit>
    {
        public override Animal CreateAnimal(int food, int number, bool isHealthy)
        {
            return new Rabbit(food, number, isHealthy, 7);
        }
    }
}