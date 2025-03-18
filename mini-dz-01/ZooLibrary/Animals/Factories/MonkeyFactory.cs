using System;

namespace ZooLibrary.Animals.Factories
{
    public class MonkeyFactory : AnimalFactoryBase<Monkey>
    {
        public override Animal CreateAnimal(int food, int number, bool isHealthy)
        {
            return new Monkey(food, number, isHealthy, 3);
        }
    }
}